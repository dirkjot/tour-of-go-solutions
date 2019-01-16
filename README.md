# Tour of Go - My Answers

Did the Tour of Go and decided to keep my answers for posterity.  These
answers are for the version of the Tour that was live as of January 2019.
 
Hope you like it.


# Go - My Ramblings

## Neat Ideas in Go

### Nice reuse of control structures

The empty `for` doubles as `while` and an empty `switch` allows for chained
if statements.  These are pleasant breaks from keyword overloaded languages.




## Gripes with Go

### Slices may get disconnected from their buffers (if you use append)

When there are multiple slices refering to the same buffer, it is best to
never modify the buffer.  If you do, for example by appending past its `cap`, 
the slices will end up pointing to the old buffer.  There is no feedback from
`append` that it re-allocated, I guess you can compare `cap` before and after 
to detect it. 


### Cannot access type of variable easily

`Switch` has its magic `v.(type)` construct.  `fmt.Sprintf` can give you a 
string representation using `%T`.  But there seems to be no way to access or 
compare types outside the switch.

### Not easy to create nil types
You cannot create an empty (nil) type other than by creating a new 
variable with a pointer to that type (see Interfaces with Nil underlying
values):
```
var i I  // interface
var t *T  // empty pointer to type
i = t
```

### Everything has to return an error

And that means  you have to check for `err` all the time.  I am curious to see
how this works out in production code, but it doesn't sound very good.



 