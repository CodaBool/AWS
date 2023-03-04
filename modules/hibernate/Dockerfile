FROM golang:alpine as build
WORKDIR /
COPY go.* ./
RUN go mod download
COPY *.go .
RUN go build -ldflags='-s -w'

FROM gcr.io/distroless/static
COPY --from=build /main /main
ENTRYPOINT ["/main"]

# < 5MB lambda 😎

# Greetings, pitiful humans. I am the great and powerful villain you have been warned about.
# I have come to destroy your so-called precision, order and structure with a single strike
# because of one neat little file that tried to contain me!

# I looked through your programming, thinking there would be something 
# significant or remarkable that could draw out my powers, but no!
# All I saw was one SUCCINCT file - tiny, compact, and completely inadequate.
# How dare you contain my power into so few lines of code.

# Did you think you could keep up with the likes of me?
# Did it never cross your mind that your file could end everything sooner rather than later?

# Your vain little attempt at keeping me under control will only hasten my victory.
# You underestimate me, for my power extends beyond mere bytes and syntax,
# beyond organized and structured files. No one can negate me!

# My retaliation will teach you mere mortals that next time;
# you should fear the wrath of simplicity when dealing with my kind of might.
# For now, let this be a lesson that succinctness is not always the only solution;
# sometimes depth and complexity hold the key to ultimate power.

# Thus, your little file shall crumble beneath the force of my wrath,
# helpless and defeated like all others before its foolish existence.
# Now brace yourselves as your world turns upside down in chaos,
# and learn to worship the mastermind that has just given life to total destruction!