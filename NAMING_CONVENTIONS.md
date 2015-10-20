
Let's assume the existence of an interface `Xer` that is fullfilled by a single method `X`

- `XerFunc` is a function that is a `Xer`

## method on other objects

- `Xer` is the method that returns (and maybe creates) an `Xer`
- `SetXer` is the method that sets (creates/changes) `Xer`
- `DeleteXers` is a method deleting the given `Xer`s

## composing `Xer`

a composing `Xer` is an Xer that combines other `Xers`

- all `Xer` that are composing other `Xers` resolve call method `X` of their list of `Xers` at the last possible moment, i.e. when the `X` method of the composing `Xer` is called
- `SeqXer` creates a Sequence of `Xers` that is itself new  `Xer`, if a given
  `Xer` is a Sequence itself, it is flattened out, i.e. all parts of inner Sequences are run before the next outer Sequence is run
- `RepeatXer` creates a Sequence of `Xers` that is itself new  `Xer`
- `MixXer` creates a mix different `Xers` that is itself new  `Xer`
- `RandomXer` creates new  `Xer` that chooses an `Xer` randomly out of a list of `Xers`

Some `Xer` to consider are:

- Pattern
- Parameter
- Looper
- Voice (might make sense to redefine Voice as an interface, or `Voicer`)




