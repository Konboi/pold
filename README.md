# pold

pold is markdown based blog tool.

# How to use.

## Init blog

```
$ pold setup

$ tree
# todo
```

## Write a article

```
$ date
Sun May  1 17:53:43 JST 2016

$ pold post new
create post/2016/05/01/175343.md

$ emacs post/2016/05/01/175343.md

$ pold server -c config.yml # default loading pold.yml
```

## Check post

Open `localhost:<port>/2016/05/01/175343` by browser.


# TODO
