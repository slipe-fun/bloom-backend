import requests
from config_tests import *
import tests
import sys
import tty
import termios

# Usage note:
# Pressing weird keys may break something, if that happens then: (sorted by severity) press a few letter or number buttons, restart the program, run the `clear` command, restart the terminal.


# Variable naming notes:
# CW means Currently Written, so this variable is written into and probably in a loop.
# CC means Currently Constructed, so it's a loop that appends many values in a list or a string or so, usually in function to later return the value

def getch():
    fd = sys.stdin.fileno()
    old_settings = termios.tcgetattr(fd)
    try:
        tty.setraw(sys.stdin)
        ch = sys.stdin.read(1)
    finally:
        termios.tcsetattr(fd, termios.TCSADRAIN, old_settings)
    return ch

_positive_signs = ("y", "e", "s", "j", "a", "+", "1", "д", "а") # Yes you can abbreviate "yes" as "s" if you're a psychopath, I guess we aren't discriminating psycopaths... You can even do the same in cyrillic, that's even better.
_negative_signs = ["n", "o", "-", "0", "н", "е", "т"]

def yrn_parse(input_text):
    cc_count = 0
    for sign in _positive_signs:
        cc_count += sign in input_text
    for sign in _negative_signs:
        cc_count -= sign in input_text
    if cc_count == 0:
        return None
    elif cc_count > 0:
        return 1
    else:
        return 0

def yrn(question): # Yes oR No
    _inp = ""
    while yrn_parse(_inp) is None:
        _inp = input(question+" [y/n]: ")
    return yrn_parse(_inp)

def to_method(text):
    if "p" in text:
        return "POST"
    elif "g" in text:
        return "GET"
    else:
        return None

def format_choices_dict(choices_dict):
    cc_text = "["
    for key in list(choices_dict.keys()):
        #ind = choices_dict[key].index(key)+1
        #left_part = choices_dict[key][:ind]
        #right_part = choices_dict[key][ind-1:]
        #cc_text += left_part+"\x1b[1m"+key.upper()+right_part+"\x1b[0m"
        cc_text += choices_dict[key]
        cc_text += "/"
    return cc_text[:-1]+"]"

def input_until_correct(text: str, choices_dict: dict):
    _choice = input(text+format_choices_dict(choices_dict)+": ")
    while _choice not in choices_dict.keys():
        _choice = input(text+format_choices_dict(choices_dict)+": ").lower()
    return _choice

def try_request(method, route):
    url = server+route

    if method == "POST":
        try:
            resp = requests.post(url)
        except requests.exceptions.ConnectionError:
            print("Make sure that the server is running.")
            exit()
    else:
        try:
            resp = requests.get(url)
        except requests.exceptions.ConnectionError:
            print("Make sure that the server is running.")
            exit()

    status_code = resp.status_code
    if status_code == 404:
        print("You probably typed the route incorrectly, the status code is 404.")
    response_text = resp.text
    return status_code, response_text

# Testing input codes
"""
for _ in range(9):
    print(repr(getch()))
exit()
"""

missing_routes = tests.t.list_missing(False)
missing_routes_text = [x[0] for x in missing_routes]

while True:
    # Autocomplete Route input
    # route = input("Route: ")
    print("\n\x1b[1mEnter the route to generate its testing code:\x1b[0m ")
    cw_route = "/"
    best_guess = missing_routes[0][0]
    print(end=f"/\x1b[2m{best_guess[1:]}\x1b[0m", flush=True)
    last_deleted = None
    while True:
        ch = getch()
        if ch == "\x03": # CTRL-C
            raise KeyboardInterrupt('Used CTRL-C')
        elif ch == "\x1b":
            ch2 = getch()
            if ch2 == "[":
                ch3 = getch()
                currently_at = missing_routes_text.index(best_guess)
                if ch3 == "A" and currently_at != 0:
                    cw_route = missing_routes[currently_at-1][0]
                elif ch3 == "B" and currently_at < len(missing_routes):
                    cw_route = missing_routes[currently_at+1][0]
                elif ch3 == "D":
                    last_deleted = cw_route
                    cw_route = "/"
                elif ch3 == "C":
                    if last_deleted is not None:
                        cw_route = last_deleted
                        last_deleted = None
                    else:
                        try:
                            cw_route += best_guess[len(cw_route)]
                        except IndexError:
                            pass
                print(end="\r"+cw_route+"\x1b[0m\x1b[J", flush=True)
        elif ch == "\r": # Enter
            break
        elif ch == "\x7f": # Backspace
            cw_route = cw_route[:-1]
        elif ch == "\t": # Tab
           cw_route = best_guess
           break
        else: # I could make it run faster but no, not like somebody pays me a salary to do this well.
            cw_route += ch
            last_deleted = None
        best_guess = sorted(missing_routes,
            key=lambda x: sum([a==b for a, b in zip(cw_route, x[0])])
            , reverse=True)[0][0]
        hint = best_guess[len(cw_route):]
        print(end="\r"+cw_route+"\x1b[2m"+hint+"\x1b[0m\x1b[J", flush=True)
    print("\n")

    method = ""
    if cw_route in missing_routes_text:
        method = sorted(missing_routes,
            key=lambda x: x==cw_route
            , reverse=True)[0][1]
    route = cw_route
    print("Editing"+" "*(method!="")+method, cw_route)
    if method == "":
        method_input = input("What is the request method? [GET/POST]: ").lower()
        _method = to_method(method_input)
        while _method is None:
            method_input = input("Wtf is that supposed to mean? [GET/POST]: ").lower()
            _method = to_method(method_input)
        method = _method
    
        print("Editing"+" "*(method!="")+method, cw_route)
    print()

    status_code, response = try_request(method, route)
    print("E - expect an error given some particular request data or none at all")
    print("W - Enter request and write the expected response manually")
    print("L - Everything is behaving correctly right now, make sure tests notice if something broke")
    behavior = input_until_correct("Desired behavior", {"e": "\x1b[1mE\x1b[0mrror", "w": "\x1b[1mW\x1b[0mrite", "l": "\x1b[1mL\x1b[0mike now"})
    if behavior == "e":
        request_data = input("Request data (just press enter if none): ")

        # Generation
        print("\n")
        if request_data.strip() == "":
            request_data = None

        if method == "GET":
            result_code = f"t.g({route}, result=\"error\", request_data={repr(request_data)}, status_code={status_code})"
        elif method == "POST":
            result_code = f"t.p({route}, result=\"error\", request_data={repr(request_data)}, status_code={status_code})"
    elif behavior == "w":
        request_data = input("Request data (just press enter if none): ")
        response_data = input("Response data: ")

        # Generation
        print("\n")
        if request_data.strip() == "":
            request_data = None

        if method == "GET":
            result_code = f"t.g({route}, result=\"error\", request_data={repr(request_data)}, status_code={status_code})"
        elif method == "POST":
            result_code = f"t.p({route}, result={repr(response_data)}, request_data={repr(request_data)}, status_code={status_code})"
    elif behavior == "l":
        request_data = input("Request data (just press enter if none): ")

        # Generation
        print("\n")
        if request_data.strip() == "":
            request_data = None
        if method == "GET":
            result_code = f"t.g({route}, result={repr(response)}, request_data={repr(request_data)}, status_code={status_code})"
        elif method == "POST":
            result_code = f"t.p({route}, result={repr(response)}, request_data={repr(request_data)}, status_code={status_code})"
        else:
            raise NameError("Invalid method, what could possibly cause this??")
    print(result_code)

    print("\n")
    if yrn("Want me to add this line to tests.py for you?"):
        with open("tests.py", "r") as file:
            tests_code = file.read()
        ANCHOR_COMMENT = "# Test generator, please add the tests here:"
        try:
            add_at = tests_code.index(ANCHOR_COMMENT)+1
        except:
            print("Idfk where to add dat shit, man... Can you comment it?")
            print(f"Just write `{ANCHOR_COMMENT}` and I'll add stuff on the next line.")
            exit()
        result_file_code = tests_code[:add_at]+result_code+"\n"+tests_code[add_at:]
        with open("tests.py", "w") as file: file.write(result_file_code)
