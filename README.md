Small commandline utility to view who is following whom on twitter

## Usage

```
followertools follows|friends|connection|commonfriends|commonfriendscount|path <user1> <user2>
```

## Examples

```
$ ./followertools follows nervengiftlabs letsencrypt
yes
$ ./followertools follows letsencrypt nervengiftlabs
no
$ ./followertools friends nervengiftlabs letsencrypt
no
$ ./followertools connection nervengiftlabs letsencrypt
☠ Nervengift (@nervengiftlabs) follows Let's Encrypt (@letsencrypt)
Let's Encrypt (@letsencrypt) doesn't follow ☠ Nervengift (@nervengiftlabs)

