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

So just found out my implementation is horribly slow for this input compared to the standard libraries implementation. I wonder why that is the
case.


So I looked at the official go package and I can't see any reason why mine should be significantly slower. I initially thought it could be that I used
recursion but they used recursion too. I am trying out benchmarkig to figure this out.


So that was interesting. I figured out why it was slow. It was because for some reason my handwritten strconv.Unquote() didn't perform anywhere close to the
stdlib's implementation. I find that a bit worrying becaue it means I haven't really learned how to write performant Go code yet. I still need to learn more about benchmarking
but boy did it work. It has also be nice to get acquainted with using a debugger.

At this point where I am really proud of what I have written. It is nice and clear without sacrificing any performance.

The other thing I was proud of what read Armin's post about stream processing and realizing that my structure is so good it wouldn't be too hard to support it as long as I could implement the
iterator interface.


Wow. I just took a peek at the VsCode debugger docs and Wow.

The case against adding AdvancePast and AdvancePastWhiteSpace - now when I declare the interface, I include two methods that are there for convenience and not part of the basis operation (by this I mean,
they can be implemented in terms of other operations)


Just realised the validate and load are actually two different usecases. Validate can be used with tools that view json to produce useful errors in a UI for example.

Also noticed my implementation is slower even though I don't do the validation step. Why? and it's also puzzling because my loading of strings was 3 times as fast. Was it fast because of the no validation?


Wow, circleci was a breeze to use :) 
