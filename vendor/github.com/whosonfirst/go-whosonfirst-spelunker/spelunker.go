package spelunker

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-roster"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

var spelunker_roster roster.Roster

// SpelunkerInitializationFunc is a function defined by individual spelunker package and used to create
// an instance of that spelunker
type SpelunkerInitializationFunc func(ctx context.Context, uri string) (Spelunker, error)

// Spelunker is an interface for writing data to multiple sources or targets.
type Spelunker interface {
	// Retrieve an individual Who's On First record by its unique ID
	GetById(context.Context, int64) ([]byte, error)
	// Retrieve an alternate geometry record for a Who's On First record by its unique ID.
	GetAlternateGeometryById(context.Context, int64, *uri.AltGeom) ([]byte, error)
	// Retrieve all the Who's On First record that are a descendant of a specific Who's On First ID.
	GetDescendants(context.Context, pagination.Options, int64, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	GetDescendantsFaceted(context.Context, int64, []Filter, []*Facet) ([]*Faceting, error)
	// Return the total number of Who's On First records that are a descendant of a specific Who's On First ID.
	CountDescendants(context.Context, int64) (int64, error)
	// Retrieve all the Who's On First records that match a search criteria.
	Search(context.Context, pagination.Options, *SearchOptions) (spr.StandardPlacesResults, pagination.Results, error)
	// Retrieve all the Who's On First records that have been modified with a window of time.
	GetRecent(context.Context, pagination.Options, time.Duration, []Filter) (spr.StandardPlacesResults, pagination.Results, error)

	GetPlacetypes(context.Context) (*Faceting, error)
	GetConcordances(context.Context) (*Faceting, error)

	HasPlacetype(context.Context, pagination.Options, *placetypes.WOFPlacetype, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	HasPlacetypeFaceted(context.Context, pagination.Options, *placetypes.WOFPlacetype, []Filter, []*Facet) ([]*Faceting, error)

	// Update this to expect *Concordance instead of parts
	HasConcordance(context.Context, pagination.Options, string, string, string, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	HasConcordanceFaceted(context.Context, pagination.Options, string, string, string, []Filter, []*Facet) ([]*Faceting, error)

	// TBD...
	// Unclear whether this should implement all of https://github.com/whosonfirst/go-whosonfirst-spatial/blob/main/spatial.go#L11
	// or https://github.com/whosonfirst/go-whosonfirst-spatial/blob/main/database/database.go#L16
	//
	// See also:
	// https://github.com/whosonfirst/go-whosonfirst-spatial-pip/blob/main/http/api/pointinpolygon.go
	// which in turns requires implementing https://github.com/whosonfirst/go-whosonfirst-spatial/blob/main/app/app.go#L21
	//
	// So it all starts to be a bit much...
	//
	// Maybe all we want are the structs and helper methods from this
	// https://github.com/whosonfirst/go-whosonfirst-spatial-pip/blob/main/pip.go
	//
	// But as twisty as all the spatial database stuff is when you start trying to make a simpler
	// version you just always end up with the same problems and questions...
	//
	// See also:
	// https://github.com/whosonfirst/go-whosonfirst-spatial-pmtiles/blob/main/database.go
	//
	// PointInPolygon(context.Context, orb.Point) (spr.StandardPlacesResults, error)

	// Not implemented yet

	/*
		GetCurrent(context.Context) ([][]byte, error)

		GetPlacetypes(context.Context) ([]placetypes.WOFPlacetype, error)
		GetPlacetype(context.Context, string) (placetypes.WOFPlacetype, error)

		GetFacetsForRecent(context.Context) (*Facets, error)
		GetFacetsForDescendants(context.Context, int64) (*Facets, error)
		GetFacetsForCurrent(context.Context) (*Facets, error)
		GetFacetsForPlacetype(context.Context) (*Facets, error)
		GetFacetsForSearch(context.Context, *SearchOptions) (*Facets, error)

		GetLanguages(context.Context) ([]*Language, error)
		GetLanguage(context.Context, string) (*Language, error)
	*/

	// TBD

	/*

		www/server.py:@app.route("/languages", methods=["GET"])
		www/server.py:@app.route("/languages/", methods=["GET"])
		www/server.py:@app.route("/languages/spoken", methods=["GET"])
		www/server.py:@app.route("/languages/spoken/", methods=["GET"])
		www/server.py:@app.route("/languages/<string:lang>", methods=["GET"])
		www/server.py:@app.route("/languages/<string:lang>/", methods=["GET"])
		www/server.py:@app.route("/languages/<string:lang>/facets", methods=["GET"])
		www/server.py:@app.route("/languages/<string:lang>/facets/", methods=["GET"])
		www/server.py:@app.route("/languages/spoken/<string:lang>", methods=["GET"])
		www/server.py:@app.route("/languages/spoken/<string:lang>/", methods=["GET"])
		www/server.py:@app.route("/languages/<string:lang>/spoken", methods=["GET"])
		www/server.py:@app.route("/languages/<string:lang>/spoken/", methods=["GET"])
		www/server.py:@app.route("/languages/spoken/<string:lang>/facets", methods=["GET"])
		www/server.py:@app.route("/languages/spoken/<string:lang>/facets/", methods=["GET"])
		www/server.py:@app.route("/languages/<string:lang>/spoken/facets", methods=["GET"])
		www/server.py:@app.route("/languages/<string:lang>/spoken/facets/", methods=["GET"])
		www/server.py:@app.route("/concordances/", methods=["GET"])
		www/server.py:@app.route("/concordances/", methods=["GET"])
		www/server.py:@app.route("/concordances/<string:who>", methods=["GET"])
		www/server.py:@app.route("/concordances/<string:who>/", methods=["GET"])
		www/server.py:@app.route("/concordances/<string:who>/facets", methods=["GET"])
		www/server.py:@app.route("/concordances/<string:who>/facets/", methods=["GET"])
		www/server.py:@app.route("/geonames/", methods=["GET"])
		www/server.py:@app.route("/gn/", methods=["GET"])
		www/server.py:@app.route("/geoplanet/", methods=["GET"])
		www/server.py:@app.route("/gp/", methods=["GET"])
		www/server.py:@app.route("/woe", methods=["GET"])
		www/server.py:@app.route("/woe/", methods=["GET"])
		www/server.py:@app.route("/tgn/", methods=["GET"])
		www/server.py:@app.route("/wikidata/", methods=["GET"])
		www/server.py:@app.route("/wd/", methods=["GET"])
		www/server.py:@app.route("/geoplanet/id/<int:id>", methods=["GET"])
		www/server.py:@app.route("/geoplanet/id/<int:id>/", methods=["GET"])
		www/server.py:@app.route("/woe/id/<int:id>", methods=["GET"])
		www/server.py:@app.route("/woe/id/<int:id>/", methods=["GET"])
		www/server.py:@app.route("/geonames/id/<int:id>", methods=["GET"])
		www/server.py:@app.route("/geonames/id/<int:id>/", methods=["GET"])
		www/server.py:@app.route("/quattroshapes/id/<int:id>", methods=["GET"])
		www/server.py:@app.route("/quattroshapes/id/<int:id>/", methods=["GET"])
		www/server.py:@app.route("/factual/id/<id>", methods=["GET"])
		www/server.py:@app.route("/factual/id/<id>/", methods=["GET"])
		www/server.py:@app.route("/simplegeo/id/<id>", methods=["GET"])
		www/server.py:@app.route("/simplegeo/id/<id>/", methods=["GET"])
		www/server.py:@app.route("/sg/id/<id>", methods=["GET"])
		www/server.py:@app.route("/sg/id/<id>/", methods=["GET"])
		www/server.py:@app.route("/faa/id/<id>", methods=["GET"])
		www/server.py:@app.route("/faa/id/<id>/", methods=["GET"])
		www/server.py:@app.route("/iata/id/<id>", methods=["GET"])
		www/server.py:@app.route("/iata/id/<id>/", methods=["GET"])
		www/server.py:@app.route("/icao/id/<id>", methods=["GET"])
		www/server.py:@app.route("/icao/id/<id>/", methods=["GET"])
		www/server.py:@app.route("/ourairports/id/<int:id>", methods=["GET"])
		www/server.py:@app.route("/ourairports/id/<int:id>/", methods=["GET"])
		www/server.py:@app.route("/id/<int:id>/descendants", methods=["GET"])
		www/server.py:@app.route("/id/<int:id>/descendants/", methods=["GET"])
		www/server.py:@app.route("/id/<int:id>/descendants/facets", methods=["GET"])
		www/server.py:@app.route("/id/<int:id>/descendants/facets/", methods=["GET"])
		www/server.py:@app.route("/megacities", methods=["GET"])
		www/server.py:@app.route("/megacities/", methods=["GET"])
		www/server.py:@app.route("/megacities/facets", methods=["GET"])
		www/server.py:@app.route("/megacities/facets/", methods=["GET"])
		www/server.py:@app.route("/nullisland", methods=["GET"])
		www/server.py:@app.route("/nullisland/", methods=["GET"])
		www/server.py:@app.route("/nullisland/facets", methods=["GET"])
		www/server.py:@app.route("/nullisland/facets/", methods=["GET"])
		www/server.py:@app.route("/placetypes", methods=["GET"])
		www/server.py:@app.route("/placetypes/", methods=["GET"])
		www/server.py:@app.route("/placetypes/<placetype>", methods=["GET"])
		www/server.py:@app.route("/placetypes/<placetype>/", methods=["GET"])
		www/server.py:@app.route("/placetypes/<placetype>/facets", methods=["GET"])
		www/server.py:@app.route("/placetypes/<placetype>/facets/", methods=["GET"])
		www/server.py:@app.route("/machinetags", methods=["GET"])
		www/server.py:@app.route("/machinetags/", methods=["GET"])
		www/server.py:@app.route("/machinetags/namespaces", methods=["GET"])
		www/server.py:@app.route("/machinetags/namespaces/", methods=["GET"])
		www/server.py:@app.route("/machinetags/namespaces/<string:ns>", methods=["GET"])
		www/server.py:@app.route("/machinetags/namespaces/<string:ns>/", methods=["GET"])
		www/server.py:@app.route("/machinetags/namespaces/<string:ns>/predicates", methods=["GET"])
		www/server.py:@app.route("/machinetags/namespaces/<string:ns>/predicates/", methods=["GET"])
		www/server.py:@app.route("/machinetags/namespaces/<string:ns>/values", methods=["GET"])
		www/server.py:@app.route("/machinetags/namespaces/<string:ns>/values/", methods=["GET"])
		www/server.py:@app.route("/machinetags/predicates", methods=["GET"])
		www/server.py:@app.route("/machinetags/predicates/", methods=["GET"])
		www/server.py:@app.route("/machinetags/predicates/<string:pred>", methods=["GET"])
		www/server.py:@app.route("/machinetags/predicates/<string:pred>/", methods=["GET"])
		www/server.py:@app.route("/machinetags/predicates/<string:pred>/namespaces", methods=["GET"])
		www/server.py:@app.route("/machinetags/predicates/<string:pred>/namespaces/", methods=["GET"])
		www/server.py:@app.route("/machinetags/predicates/<string:pred>/values", methods=["GET"])
		www/server.py:@app.route("/machinetags/predicates/<string:pred>/values/", methods=["GET"])
		www/server.py:@app.route("/machinetags/values", methods=["GET"])
		www/server.py:@app.route("/machinetags/values/", methods=["GET"])
		www/server.py:@app.route("/machinetags/values/<string:value>", methods=["GET"])
		www/server.py:@app.route("/machinetags/values/<string:value>/", methods=["GET"])
		www/server.py:@app.route("/machinetags/values/<string:value>/namespaces", methods=["GET"])
		www/server.py:@app.route("/machinetags/values/<string:value>/namespaces/", methods=["GET"])
		www/server.py:@app.route("/machinetags/values/<string:value>/predicates", methods=["GET"])
		www/server.py:@app.route("/machinetags/values/<string:value>/predicates/", methods=["GET"])
		www/server.py:@app.route("/machinetags/places/<string:ns_or_mt>", methods=["GET"])
		www/server.py:@app.route("/machinetags/places/<string:ns_or_mt>/", methods=["GET"])
		www/server.py:@app.route("/machinetags/places/<string:ns>/<string:pred>", methods=["GET"])
		www/server.py:@app.route("/machinetags/places/<string:ns>/<string:pred>/", methods=["GET"])
		www/server.py:@app.route("/machinetags/places/<string:ns>/<string:pred>/<string:value>", methods=["GET"])
		www/server.py:@app.route("/machinetags/places/<string:ns>/<string:pred>/<string:value>/", methods=["GET"])
		www/server.py:@app.route("/tags", methods=["GET"])
		www/server.py:@app.route("/tags/", methods=["GET"])
		www/server.py:@app.route("/names", methods=["GET"])
		www/server.py:@app.route("/names/", methods=["GET"])
		www/server.py:@app.route("/tags/<tag>", methods=["GET"])
		www/server.py:@app.route("/tags/<tag>/", methods=["GET"])
		www/server.py:@app.route("/tags/<tag>/fatets", methods=["GET"])
		www/server.py:@app.route("/tags/<tag>/facets/", methods=["GET"])
		www/server.py:@app.route("/categories/<category>", methods=["GET"])
		www/server.py:@app.route("/categories/<category>/", methods=["GET"])
		www/server.py:@app.route("/categories/<category>/facets", methods=["GET"])
		www/server.py:@app.route("/categories/<category>/facets/", methods=["GET"])
		www/server.py:@app.route("/postalcode/<code>", methods=["GET"])
		www/server.py:@app.route("/postalcode/<code>/", methods=["GET"])
		www/server.py:@app.route("/postalcodes/<code>", methods=["GET"])
		www/server.py:@app.route("/postalcodes/<code>/", methods=["GET"])
		www/server.py:@app.route("/postalcode/<code>/facets", methods=["GET"])
		www/server.py:@app.route("/postalcode/<code>/facets/", methods=["GET"])
		www/server.py:@app.route("/postalcodes/<code>/facets", methods=["GET"])
		www/server.py:@app.route("/postalcodes/<code>/facets/", methods=["GET"])
		www/server.py:@app.route("/opensearch", methods=["GET"])
		www/server.py:@app.route("/opensearch/", methods=["GET"])
		www/server.py:@app.route("/opensearch/<scope>", methods=["GET"])
		www/server.py:@app.route("/opensearch/<scope>/", methods=["GET"])
		www/server.py:@app.route("/search", methods=["GET"])
		www/server.py:@app.route("/search/", methods=["GET"])
		www/server.py:@app.route("/auth", methods=["GET"])
		www/server.py:@app.route("/auth/", methods=["GET"])
		www/server.py:@app.route("/search/facets", methods=["GET"])
		www/server.py:@app.route("/search/facets/", methods=["GET"])

	*/
}

// RegisterSpelunker registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Spelunker` instances by the `NewSpelunker` method.
func RegisterSpelunker(ctx context.Context, scheme string, init_func SpelunkerInitializationFunc) error {

	err := ensureSpelunkerRoster()

	if err != nil {
		return err
	}

	return spelunker_roster.Register(ctx, scheme, init_func)
}

func ensureSpelunkerRoster() error {

	if spelunker_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		spelunker_roster = r
	}

	return nil
}

// NewSpelunker returns a new `Spelunker` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `SpelunkerInitializationFunc`
// function used to instantiate the new `Spelunker`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterSpelunker` method.
func NewSpelunker(ctx context.Context, uri string) (Spelunker, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := spelunker_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(SpelunkerInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpelunkerRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range spelunker_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
