Slice
=====

Slice is a framework that makes it easy to build applications out
of reusable, composable modules.

## Application lifecycle

1. init slice configuration: debug, env, name, description, etc.
2. provide default variables: Name, Debug, Env, Log, Args
3. provide default components: context
4. load bundles configuration from environment
5. build bundles: provide bundle types
6. boot bundles: invoke custom bundle functions
7. invoke dispatch function
8. shutdown bundles: invoke custom bundle functions
9. cleanup container

## Variables

Slice contains variables provided by default.

#### `slice.Name`

The name of application.

#### `slice.Debug`

The debug flag. Will be loaded from environment variable `DEBUG`.

#### `slice.Env`

The run environment. May be `dev`, `test`, `prod`. Default: `prod`. Will
be loaded from environment variable `ENV`.

#### `slice.Log`

A slice system log function. For `dev` environment logging purposes. Must
be used only in bundle scope.

#### `slice.Args`

An arguments of application. In common alias for `os.Args`.