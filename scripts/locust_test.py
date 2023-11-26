# from dataclasses import dataclass, field
# import queue
import random
from locust import HttpLocust
from locust import TaskSet
from locust import task
from locust import events
from locust.events import EventHook
# from locust import between
import json

# For HTML reporting
from locust.web import app
import urllib
from src import report
app.add_url_rule('/htmlreport', 'htmlreport', report.download_report)

def custom_request_failure_handler(request_type, name, response_time, exception):
    print(f"Request Failure: {request_type} {name} {response_time} {exception}")
    if exception.response.status_code >= 400 and exception.response.status_code < 500:
        events.request_success.fire(
            request_type=request_type,
            name=name,
            response_time=response_time,
            response_length=0,
        )

events.request_failure = EventHook()
events.request_failure += custom_request_failure_handler

class TestInput:
    people_list: list = []
    dropoff_list: list = []
    cars: list = []

    def __init__(self):
        # self.people_list = self.generate_people(100000)
        # self.cars = self.generate_cars(10000)
        self.people_list = self.generate_people(1000)
        self.cars = self.generate_cars(400)
        
    def delay(self, seconds) -> None:
        for i in range(seconds):
            print(i)
            self.wait(1000)

    def generate_cars(self, number_of_cars) -> list:
        return [
            {"id": i, "seats": random.randint(4, 6)}
            for i in range(1, number_of_cars + 1)
        ]

    def generate_people(self, number_of_people) -> list:
        return [
            {"id": i, "people": random.randint(1, 6)}
            for i in range(1, number_of_people + 1)
        ]
    
    def pop_waiting_people(self) -> dict:
        if len(self.people_list) > 0:
            item = self.people_list.pop(0)
            return item
        else:
            return None
        
    def pop_dropoff_people(self) -> dict:
        if len(self.dropoff_list) > 0:
            item = self.dropoff_list.pop(0)
            return item
        else:
            return None
        

global_input : TestInput
run_once = False

class CarsTaskSet(TaskSet):

    def cars_request(self):
        global global_input

        print(f"/cars {json.dumps(global_input.cars)}")

        response = self.client.put(
            "/cars",
            json=global_input.cars
        )
        
        print(f"/cars {response.status_code}")


class CarPoolingTaskSet(TaskSet):

    @task(1)
    def journey_request(self):
        global global_input
        people = global_input.pop_waiting_people()

        if people:
            print(f"/journey pop people {people['id']}")
        else:
            print(f"/journey pop None")

        if people:
            response = self.client.post(
                "/journey",
                json=people
            )
            
            print(f"/journey {response.status_code}")
            
            if response.status_code in [200, 202, 204]:
                global_input.dropoff_list.append(people)
            else:
                global_input.people_list.append(people)
            

    @task(1)
    def dropoff_request(self):
        global global_input
        people = global_input.pop_dropoff_people()

        if people:
            print(f"/dropoff pop people {people['id']}")
        else:
            print(f"/dropoff pop None")

        if people:
            data = {"ID": people["id"]}
            form_data = urllib.parse.urlencode(data)

            headers = {"Content-Type": "application/x-www-form-urlencoded"}
            response = self.client.post(
                "/dropoff",
                data=form_data,
                headers=headers
            )

            print(f"/dropoff {response.status_code}")
            
            if response.status_code in [200, 202, 204]:
                global_input.people_list.append(people)
            else:
                global_input.dropoff_list.append(people)


class MyLocust(HttpLocust):
    task_set = CarPoolingTaskSet
    min_wait = 1000
    max_wait = 1000

    def __init__(self):
        super().__init__()
        global run_once
        global global_input
        if not run_once:
            global_input = TestInput()
            cars_task_set = CarsTaskSet(self)
            cars_task_set.cars_request()
            run_once = True
