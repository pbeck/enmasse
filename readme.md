# enmasse

`enmasse` (/ɑn ˈmæs/) is a utility application for creating GMail draft emails compiled from [golang templates](https://golang.org/pkg/text/template/) and JSON data. Read more about enmasse at [beckman.io/enmasse](http://beckman.io/enmasse)

**ProTip:** use [Boomerang for GMail](http://www.boomeranggmail.com/) for scheduling emails.

Created by Pelle Beckman, [@pbeck](http://twitter.com/pbeck)

***Please don’t use enmasse for sending spam!***

## Installation

You need to generate your own Google Gmail API Credentials – [there’s a great tutorial available at Google Developers](https://developers.google.com/gmail/api/quickstart/go).

Download binaries or build from source. Keep the `client_secret.json` in the same directory as binary. On the first run enmasse will ask you to open a browser and perform authentication.

## Example usage

`enmasse -template=template.tmpl -data=addresses.json`

**addresses.json:**

    [{
	    "first_name": "Alistair",
	    "email": "alistair.hennessey@gmail.com",
	    "title": "PhD"
    }, {
    	"first_name": "Klaus",
	    "email": "klaus.daimler@gmail.com",
	    "title": "Research Assistant",
	    "label": "Team Zissou"
    }, {
	    "first_name": "Ned",
	    "email": "ned.plimpton@gmail.com",
	    "location": "Port-au-Patois",
	    "label": "Team Zissou"
    }]

**template.txt:**

    Hey {{ .first_name }},

    {{ if .title }}
    I guess I should be pretty impressed that you’ve achieved the rank of {{ .title }}...
    {{ end }}

    Now if you'll excuse me, I'm going to go on an overnight drunk,
    and in 10 days I'm going to set out to find the shark that ate my friend
    and destroy it. Anyone who wants to join me is more than welcome.

    Best regards,

    Steve
    
    P.S. I’ve attached your party invitation.

    {{ if .location }}
    P.S. 2 {{ .location}} is pretty nice this time of the year, right?
    {{ end }}
    
## Flags

    -template=FILE    Template file
    -data=FILE        JSON data file

## Scheduling messages

A built-in scheduling function would be nice, but **a)** I feel it’s out of scope, **b)** it would require a lot of work and most likely require a backend of some sort, and  **c)** [Boomerang for GMail](http://www.boomeranggmail.com/) seems to work like a charm.

## Contributing

Please report bugs or feature requests in the [issue tracker at GitHub](https://github.com/pbeck/enmasse/issues).

## License

The MIT License (MIT)

Copyright (c) 2016 Pelle Beckman, http://beckman.io

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
