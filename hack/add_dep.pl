#!/usr/bin/perl -w
use strict;

## hasdep will check if given pom.xml (string) has given artifact configured
## as a dependency.
sub hasdep {
  my $pom = shift;
  my $art = shift;
  my $dep = $pom;
  $dep =~ s#.*<dependencies>##is;
  $dep =~ s#</dependencies>.*##is;
  return ($dep =~ m/$art/is);
}

## getdep will create a dependency xml block for given groupid, artifact and
## version.
sub getdep {
  my $groupid = shift;
  my $artifact= shift;
  my $version = shift;
  my $dep = "<dependency>";
  $dep .= "<groupId>$groupid</groupId>";
  $dep .= "<artifactId>$artifact</artifactId>";
  $dep .= "<version>$version</version>";
  $dep .= "<scope>provided</scope>";
  $dep .= "</dependency>";
  return $dep;
}

## getpluginv will return the version of the configured maven kie plugin, if
## present.
sub getpluginv {
  my $pom = shift;
  my $ver = undef;
  $pom =~ s#.*<plugins>##is;
  $pom =~ s#</plugins>.*##is;
  $pom =~ s#<plugin>#\t#gis;
  foreach my $p (split(/\t/,$pom)) {
    if ($p =~ m/kie-maven-plugin/) {
      $ver = $p;
      $ver =~ s#.*<version>##is;
      $ver =~ s#<version>.*##is;
      last;
    }
  }
  return $ver;
}

## main will move to the batmobile and go.
sub main {
  my $pomfile = shift || $ENV{DEP_POMFILE};
  my $groupid = shift || $ENV{DEP_KIE_GROUP_ID} || "org.kie";
  my $artifact= shift || $ENV{DEP_KIE_ARTIFACT} || "kie-api";

  if (!$pomfile) {
    print STDERR "Warning: no pom file specified, not adding kie api dependency\n";
    exit(0);
  }

  open(F,"<".$pomfile) or die "Failed opening: $pomfile: $!";
  my $pom = join('',<F>);
  close(F);

  if ($pom =~ m/<dependencies>/is) {
    print STDERR "Info: dependencies already available\n";
  } else {
    print STDERR "Info: adding dependencies tag\n";
    $pom =~ s#<repositories>#<dependencies></dependencies><repositories>#is;
  }

  if (hasdep($pom,$artifact)) {
    print STDERR "Info: kie-maven-plugin dependency already available, nothing to do...\n";
    exit(0);
  }

  my $version = getpluginv($pom) || "6.5.0.Final-redhat-2";

  my $dep = getdep($groupid,$artifact,$version);
  $pom =~ s#</dependencies>#$dep</dependencies>#is;

  open(F,">".$pomfile) or die "Failed writing: $pomfile: $!";
  print F $pom;
  close(F);

  print "created pom.xml\n";
  print "-" x 64;#
  print "\n".$pom."\n";
}

main(shift,shift,shift,shift);
