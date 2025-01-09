# reaperstream

This software is provided for educational purposes only. Please don't use this for any malicious purposes.

# what is it

See my video discussing and demonstrating this: https://youtu.be/pgaUoboK0XY

Recently it has been discovered that law enforcement automatic license plate readers have been deployed with publicly accessible endpoints. This software is a proof of concept to show how easy it is to scrape these endpoints for data.

# where are these readers?

Based on the license plated visible in images, we can determine that we have found these cameras in Florida and Illinois.

# what does law enforcement do with this data?

Reading product literature from the vendor, we have determined that this data is collected and stored for an undetermined length of time.
The vendor appears to provide aggregated searches across all devices regardless of where they are deployed, and allows
for searching ly license plate, type of vehicle, and more. They also provide a "hot list" which can be used to alert law
enforcement when a vehicle of interest is detected.

Other features include the ability to search for vehicles that appear to be following a specific vehicle.

# why is this available to the public?

Because law enforcement is not properly securing their endpoints. A combination of laziness, stupidity, and lack of care for the privacy
of the citizens they are supposed to protect has led to this situation. If this were a private company, they would be sued into oblivion.

# but 1984 is just a book not an instruction manual

Tell me about it, sister.

# but I can see so many bad things that could be done with this data

I know, right? I won't even go into the possibilities here. I'm sure you can think of a few yourself.

# ok so how do I build this?

You need the Go compiler, if you're on Linux you can probably just do `sudo apt-get install golang` or `sudo yum install golang` or whatever. If you're on Windows or Mac you can download it from the Go website.
Use git to clone this repository, or you can download it as a zip and extract it.
Once you have Go installed, you can just run `go build` in the directory where you cloned this repository. You will now have a binary called `reaperstream` that you can run.

```
Usage of ./reaperstream:
  -debug
    	enable debug mode
  -hosts string
    	comma separated list of hosts to connect to
  -listen string
    	address:port to listen on (default "localhost:8888")
```

# how do I run this?

Watch Matt Brown's video to find camera aggregators to run this against. It took me 10 seconds to find a list, and a quick nmap found even more.
