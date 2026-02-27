import customlib_testing as t
import sys

show_not_added = True
if "quiet" in sys.argv:
    t.show_success = False
    show_not_added = False

t.g("/ws")

if show_not_added:
    t.list_missing()

t.run_tests()
