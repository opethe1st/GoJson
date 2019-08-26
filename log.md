## Learnings

I learned that Go's "=" is not good enough for asserting equality of nested objects when writing tests. This was disappointing because it meant I needed to
use a third party library.
I learned that you need to be careful about using ":=" in a loop if it shadows an outer variable. I lost a couple of hours to trying to debug this

P.s I skimmed through the official implementation and it feels like I could do a better job (Is it hard to follow or it is just me?)
Plus lots of state. My implementation probably suffers from using recursion though


### Random thoughts
Not use recursion? - is that possible?

After looking at the code, I realize that instead of passing (s, current) into virtually every function, I could create a scanner type
and have these functions be methods in this type.
```golang
type Scanner struct{
    s string
    current int
}

func (s *Scanner) consumeWhiteSpace(){

}

func (s *Scanner) isMappingStart(){

}

func (s *scanner) isMappingEnd(){

}
```
Other thing is to introduce the concept of tokens - only needed for strings? Since that's where you can have escaped characters!

So initially, I introduced the concept of scanners but then I realized that it didn't fit completely.
So I discovered the concept that was missing was the concept of iterators. A iterator is simply - struct{s string, offset int}
And this made other stuff make sense since I moved the stuff that didn't make sense in the iterator to standalone functions that instead of (s, current)
made use of the iterator. The other interesting thing is that I violated Command Query Responsibility Segregation (CQRS) and it was the right decision.
Know when to break the rules!


Probably need to make sure this handles unicode

Just realized that I didn't check to see what the form is when a file is read? Are the characters interpreted literally? Making the need for string parsing moot?
