# convertdir

`convertdir` converts a directory of images to one with Kindle or PDF formatted images.

```
$ go get github.com/grokify/mogo/image/apps/convertdir
$ which convertdir
$ convertdir -i {inputdir} -o {outputdir} -f {format (kindle|pdf)}
$ convertdir -i images/orig -o images/kindle -f kindle
```