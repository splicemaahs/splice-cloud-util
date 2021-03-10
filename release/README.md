# RELEASE Process

1. export RELEASE_VERSION=v<newSemVer>
  - `export RELEASE_VERSION=v0.0.24`
2. Create a new branch for release.
  - `git checkout -b RELEASE_${RELEASE_VERSION}`
3. `make changelog`
4. Review changelogs/releases/<newSemVer>.md
  - `ls changelog/fragments`
  - `mdcat changelog/releases/${RELEASE_VERSION}.md`
5. Create and Merge PR
  - `git add -A; git commit -m "RELEASE of ${RELEASE_VERSION}"; git push origin RELEASE_${RELEASE_VERSION}`
  - https://github.com/splicemaahs/splice-cloud-util/ Create a PR and Merge
6. Pull main
  - `git checkout main; git fetch; git pull`
7. git tag <newSemVer>
  - `git tag v0.0.24`
8. git push origin <newSemVer>
  - `git push origin v0.0.24`
9. Check to see if Release GitHub Action kicks off
10.  Check Releases Page, and homebrew-tap/Formula/splice-cloud-util.rb

