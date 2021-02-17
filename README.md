# GoSkew

GoSkew is a program for post-processing G-code to account for axis skew. The
method and techniques are heavily based on
[gSkewer](https://github.com/MechanizedMedic/gskewer), a python script with a
similar purpose.

I have personally found the method of measuring the error outlined by
[gSkewer](https://github.com/MechanizedMedic/gskewer) quite fiddly and prone to
inaccuracy. For that reason, at least for XY skew, I am working on a new method
of calculating error that requires only a printed object and calipers.

# TODO

- Add docopt's
- Add skew code
- Add and test triangle error measurement
- Write an awesome readme
