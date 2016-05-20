# pold

pold is markdown based blog tool.

under development.

# How to use.

## Init blog

```
$ pold setup

$ tree
.
├── pold.yml
├── post
└── templates
    ├── archive.html
    ├── base.html
    ├── index.html
    └── post.html
2 directories, 5 files
```

## Write a article

```
$ date
Sun May  1 17:53:43 JST 2016

$ pold new sample
create post/2016/05/01/sample.md

$ emacs post/2016/05/01/sample.md

$ pold server -c config.yml # default loading pold.yml
```

## Check post

Open `localhost:<port>/2016/05/01/175343` by browser.


# TODO

- [ ] draft mode
- [ ] tag page
