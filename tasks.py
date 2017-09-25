#!/usr/bin/env python3

from invoke import task
from pprint import pprint

import os

@task
def deploy_lambdas(ctx, state="test"):
  output = ctx.run("ls functions")
  for lambda_name in output.stdout.split("\n"):
    if len(lambda_name) == 0:
      continue
    script_dir = "{0}/functions/{1}".format(os.getcwd(), lambda_name)
    with ctx.cd(script_dir):
      if os.path.exists("{0}.zip".format(lambda_name)):
        ctx.run("mv {0}.zip {0}.zip.backup".format(lambda_name))
      venv_path = "{0}/venv".format(script_dir)
      if os.path.exists(venv_path):
        ctx.run("rm -rf {0}".format(venv_path))
      ctx.run("python3.6 -m venv venv")
      with ctx.prefix("source venv/bin/activate"):
        ctx.run("pip install -r requirements.txt")
      ctx.run("zip -r9 {0}/{1}.zip {0}/venv/lib/python3.6/site-packages/*".format(script_dir, lambda_name))
      ctx.run("zip -r9 {0}/{1}.zip {0}/{1}.py".format(script_dir, lambda_name))
      ctx.run("terraform init")
      if state == "test":
        ctx.run("terraform plan")
      elif state == "deploy":
        ctx.run("terraform apply")



@task
def deploy(ctx):
    deploy_lambdas(ctx, state="deploy")


@task
def test(ctx):
    deploy_lambdas(ctx, state="test")
