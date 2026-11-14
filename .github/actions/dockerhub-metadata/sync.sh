#!/usr/bin/env bash
# Set a Docker Hub repository's short description, overview and categories.
# Reads DH_REPOSITORY, DH_USERNAME, DH_TOKEN, DH_DESCRIPTION, and optionally
# DH_README (with DH_LINK_BASE for its relative links), DH_OVERVIEW and
# DH_CATEGORIES. Never prints the token.
set -euo pipefail
export LC_ALL=C.UTF-8

api=https://hub.docker.com/v2
die() {
	echo "dockerhub-metadata: $*" >&2
	exit 1
}

: "${DH_REPOSITORY:?}" "${DH_USERNAME:?}" "${DH_TOKEN:?}" "${DH_DESCRIPTION:?}"
[ "${#DH_DESCRIPTION}" -le 100 ] ||
	die "description has ${#DH_DESCRIPTION} characters; Docker Hub allows 100"

if [ -n "${DH_README:-}" ]; then
	[ -f "$DH_README" ] || die "no such readme: $DH_README"
	# Point relative links (not URLs, anchors or absolute paths) at GitHub.
	dir=$(dirname "$DH_README")
	base=${DH_LINK_BASE:?}
	[ "$dir" = . ] || base=$base/${dir#./}
	overview=$(sed -E "s|\]\(([^):/#][^):]*)\)|](${base}/\1)|g" "$DH_README")
else
	overview=${DH_OVERVIEW:-}
fi
[ "${#overview}" -le 25000 ] ||
	die "overview has ${#overview} characters; Docker Hub allows 25000"

# Category slugs to the {name, slug} objects the API expects.
all=$(curl -fsS "$api/categories")
categories=$(jq -c --arg want "${DH_CATEGORIES:-}" '
	($want | split(" ") | map(select(. != ""))) as $w
	| [.[] | select(.slug | IN($w[]))] as $found
	| if ($w | length) > 3 then error("at most 3 categories, got \($w | length)")
	  elif ($found | length) < ($w | length) then
	    error("unknown categories: \($w - ($found | map(.slug)) | join(" ")); valid: \(map(.slug) | join(" "))")
	  else $found | map({name, slug}) end' <<<"$all") || die "bad categories"

jwt=$(jq -n --arg u "$DH_USERNAME" --arg p "$DH_TOKEN" '{identifier: $u, secret: $p}' |
	curl -fsS -H 'Content-Type: application/json' -d @- "$api/auth/token" |
	jq -er .access_token) || die "Docker Hub login failed"

patch() {
	local code
	code=$(curl -sS -o /tmp/dockerhub-metadata.out -w '%{http_code}' -X PATCH \
		-H "Authorization: Bearer $jwt" -H 'Content-Type: application/json' \
		-d @- "$api/repositories/$DH_REPOSITORY/$1")
	[ "$code" = 200 ] || die "PATCH $1 answered $code: $(head -c 300 /tmp/dockerhub-metadata.out)"
}

jq -n --arg d "$DH_DESCRIPTION" --arg f "$overview" '{description: $d, full_description: $f}' | patch ""
patch categories <<<"$categories"

echo "dockerhub-metadata: updated $DH_REPOSITORY (categories: $(jq -r 'map(.slug) | join(" ")' <<<"$categories"))"
