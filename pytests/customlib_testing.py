import requests
from config_tests import *

# Could use a lib like unittests but not this time

show_success = True

def printfail(fail_text: str, route: str, fail_type: str = "UNKNOWN"):
    fail_left = "\x1b[31;1m"+fail_type
    spaces_count = max([10-len(fail_type), 0])
    fail_left += " "*spaces_count

    fail_left += " FAIL: "
    print(fail_left+route+"\x1b[0m")
    print("\x1b[31;1m"+fail_text+"\x1b[0m")

def printsuccess(route: str, success_type: str = ""):
    spaces_count = max([11-len(success_type), 0]) # 11 because success is formatted slightly differently and there's a space in "status code"
    if show_success:
        print(f"\x1b[32;1mCorrect {success_type}"+" "*spaces_count, f"at {route}\x1b[0m")

class Test:
    def __init__(self, route: str, result: str
            , request_data: dict = None, 
            method: str = "GET", s_code: int = None):
        """
            route: Example - "/metrics"
            result: What should the test return
            s_code: What HTTP status should be returned for the test to pass
            request_data: What data does the request also provide
        """
        self.route = route
        self.result = result
        self.request_data = request_data
        self.method = method
        self.status_code = s_code
    def __call__(self):
        if   self.method == "GET": # Could use match case
            resp = requests.get(server+self.route, self.request_data)
        elif self.method == "POST":
            resp = requests.post(server+self.route, self.request_data)
        else:
            raise TypeError("Unsupported HTTP method") # TODO Add a new exception for everything

        # Checking the status code
        should_error = "error" in self.result # Different from the response data check because this is a check for the case when the correct status code isn't set, so it's just guessing. The response data is a perfect 1-2-1 check to a json unless the user inputs the magic word "error" instead of the json to change the behavior.
        valid_response_code = False
        if self.status_code is None:
            if should_error and resp.status_code == 200:
                printfail("Expected error but s-code is 200", self.route, "STATUSCODE")
            else:
                printsuccess(self.route, "status code")
                valid_response_code = True
        else:
            if self.status_code != resp.status_code:
                printfail(f"Wrong status code {resp.status_code}, should be", self.status_code, "at" , self.route, "STATUSCODE")
            else:
                printsuccess(self.route, "status code")
                valid_response_code = True

        # Checking the response data

        response_error = "error" in resp.text
        valid_response = False
        if self.result == "error" or self.result == "err":
            valid_response = "error" in resp.text
            if valid_response:
                printsuccess(self.route, "response")
        else:
            valid_response = response.text == self.result
        return valid_response_code and valid_response


def g(route: str, result="error", request_data: dict = None, status_code: int = None): # GET request
    tests.append(Test(route=route, result=result, request_data=request_data, s_code=status_code, method="GET"))
    test_routes.update((route,))

def p(route: str, result="error", request_data: dict = None, status_code: int = None): # POST request
    tests.append(Test(route=route, result=result, request_data=request_data, s_code=status_code, method="POST"))
    test_routes.update((route,))

def list_missing(quiet: bool = False):
    with open("cmd/api/main.go") as file:
        go_code = file.read()

    start = go_code.index("authMiddleware := middleware.NewAuthMiddleware(sessionApp)")+58
    end = go_code.index("log.Fatal(fiberApp.Listen(fmt.Sprintf(\"%s:%d\", cfg.Server.Host, cfg.Server.Port)))\n}") # len: 85, but why would you need that?
    go_code = go_code[start:end]

    go_lines = go_code.split("\n")
    cc_missing_tests = [] # CurrentlyConstructed
    cc_methods = []
    for line in go_lines:
        if len(line.strip()) == 0: continue

        start_route = line.index('t("/')
        end_route = line.index('", ')
        route = (line[start_route+3:end_route])
        if route not in test_routes and not quiet:
            print("Not added:", route)
        if line[start_route-4:start_route+1]==".Post":
            method = "POST"
        elif line[start_route-4:start_route+1]=="p.Get":
            method = "GET"
        else:
            print(line)
            raise SyntaxError("What happened here?? Can't detect the method")
        
        cc_missing_tests.append((route, method))
        
    return cc_missing_tests

def run_tests():
    [test() for test in tests]

tests = []
test_routes = set()
