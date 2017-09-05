#!/usr/bin/env python3

from invoke import task
from pprint import pprint

import os

@task
def deploy_lambdas(ctx):
  output = ctx.run("ls functions")
  for lambda_name in output.stdout.split("\n"):
    script_dir = "{0}/functions/{1}".format(os.getcwd(), lambda_name)
    with ctx.cd(script_dir):
      if os.path.exists("{0}.zip".format(lambda_name)):
        ctx.run("mv {0}.zip {0}.zip.backup".format(lambda_name))
    if not os.path.exists("{0}/venv".format(script_dir)):
      ctx.run("python3 -m venv venv")
      ctx.run("source venv/bin/activate && pip install -r requirements.txt")
    with ctx.cd("cd {0}/venv/lib/python3.6/site-packages".format(script_dir)):
      ctx.run("zip -r9 {0}/{1}.zip *".format(script_dir, lambda_name))
    with ctx.cd(script_dir):
      ctx.run("zip -r9 {0}/{1}.zip {1}.py".format(script_dir, lambda_name))
