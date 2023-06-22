using Go = import "/go.capnp";
@0x93814030c230a561;
$Go.package("capnp_go");
$Go.import("blumer-ms-refers");

struct ProfileInfo {
  userId @0 :Text;
  username @1 :Text;
  isActive @2 :Bool;
  reward @3 :Float64;
}

struct ProfileWallet {
    userId @0 :Text;
}