
from rjtools.util.testing import run_modules

COMMAND_PREFIX = [
    "./cfgo",
]

def run():
    from . import parse
    from . import cfg

    return run_modules("cfgo", locals(), COMMAND_PREFIX)

