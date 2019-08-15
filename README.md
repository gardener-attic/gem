gem
===

`gem` is the *G*ardener *E*xtension *M*anager - not to be confused with
Rubygems of course.

When working with Gardener, you'll always want to use extensions. Extensions
come in the form of of `controller-registration.yaml` files. These files
register an extension provider in Gardener so that it can then make use of
its functionality when creating a Shoot cluster. More on this topic can be
found in the [Gardener repository](https://github.com/gardener/gardener/blob/master/docs/extensions/controllerregistration.md).

One issue when working with extensions is to get the right set of versions.
`gem` is a command line tool & library that addresses this issue by providing
means to define extension requirements, resolve them to exact revisions and
then download them into one `controller-registrations.yaml` file (one big
YAML that contains all `controller-registration.yaml` files as a list).

Installation
------------

To install `gem`, you need a working `go` installation `>1.11`. Once you have
it, `cd` to an empty, temporary directory (in order not to mess up any `go`
project) and run

```bash
GO111MODULE=on go get github.com/gardener/gem
```

Usage
-----

In order to use `gem` you need to define a `requirements.yaml`. An example can
be found [here](example/requirements.yaml). In this file, you list your
required gardener extensions.

For a requirement, you always have to specify a `name` and either a `revision`,
`version` or `branch`. Optionally, you can also specify a `filename`, if the
name of controller-registration is not the default
`controller-registration.yaml`.

Once you've successfully defined a `requirements.yaml` file, `gem` provides the
following commands to work with it, though the most important one will probably
be *`ensure`*:

* *`solve`*: tries to solve the requirements given in your
  `requirements.yaml` and write it into a `locks.yaml`.

* *`fetch`*: requires a `locks.yaml` to be present. Fetches the
  `controller-registrations` specified via the requirements and locks.

* *`ensure`*: ensures that the controller-registrations you specified
  in your `requirements.yaml` are present and up to date. It can optionally also
  update your dependencies to the latest allowed version.
