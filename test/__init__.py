
from rjtools.util.testing import run_modules

COMMAND_PREFIX = [
    "./golze",
]

def run():
    from . import parse
    from . import cfg

    return run_modules("golze", locals(), COMMAND_PREFIX)

