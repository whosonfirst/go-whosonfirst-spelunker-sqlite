var whosonfirst = whosonfirst || {};
whosonfirst.spelunker = whosonfirst.spelunker || {};

whosonfirst.spelunker.leaflet = (function(){

    var self = {
	'draw_point': function(map, geojson, layer_args){
	    
	    var layer = L.geoJson(geojson, layer_args);
	    
	    layer.addTo(map);
	    return layer;
	},
	
	'draw_poly': function(map, geojson, layer_args){
	    
	    var layer = L.geoJson(geojson, layer_args);
	    
	    // https://github.com/Leaflet/Leaflet.label
	    
	    try {
		var props = geojson['properties'];
		
		if (props){
		    var label = props['lflt:label_text'];
		    
		    if (label){

			var label_args = {
			    noHide: true,
			    pane: layer_args.pane,
			}
			
			layer.bindTooltip(label, label_args);			
		    }
		}
		
		else {
		    console.log("polygon is missing a properties dictionary");
		}
	    }
	    
	    catch (e){
		console.log("failed to bind label because " + e);
	    }
	    
	    layer.addTo(map);
	    return layer;
	},
	
	'draw_bbox': function(map, geojson, layer_args){
	    
	    var bbox = whosonfirst.spelunker.geojson.derive_bbox(geojson);
	    
	    if (! bbox){
		console.log("no bounding box");
		return false;
	    }
	    
	    var bbox = geojson['bbox'];
	    var swlat = bbox[1];
	    var swlon = bbox[0];
	    var nelat = bbox[3];
	    var nelon = bbox[2];
	    
	    var geom = {
		'type': 'Polygon',
		'coordinates': [[
		    [swlon, swlat],
		    [swlon, nelat],
		    [nelon, nelat],
		    [nelon, swlat],
		    [swlon, swlat],
		]]
	    };
	    
	    var bbox_geojson = {
		'type': 'Feature',
		'geometry': geom
	    }
	    
	    return self.draw_poly(map, bbox_geojson, layer_args);
	},
	
	'fit_map': function(map, geojson, force){
	    
	    var bbox = whosonfirst.spelunker.geojson.derive_bbox(geojson);
	    
	    if (! bbox){
		console.log("no bounding box");
		return false;
	    }
	    
	    if ((bbox[1] == bbox[3]) && (bbox[2] == bbox[4])){
		map.setView([bbox[1], bbox[0]], 14);
		return;
	    }
	    
	    var sw = [bbox[1], bbox[0]];
	    var ne = [bbox[3], bbox[2]];
	    
	    var bounds = new L.LatLngBounds([sw, ne]);
	    var current = map.getBounds();
	    
	    var redraw = true;
	    
	    if (! force){
		
		var redraw = false;
		
		/*
		   console.log("south bbox: " + bounds.getSouth() + " current: " + current.getSouth().toFixed(6));
		   console.log("west bbox: " + bounds.getWest() + " current: " + current.getWest().toFixed(6));
		   console.log("north bbox: " + bounds.getNorth() + " current: " + current.getNorth().toFixed(6));
		   console.log("east bbox: " + bounds.getEast() + " current: " + current.getEast().toFixed(6));
		 */
		
		if (bounds.getSouth() <= current.getSouth().toFixed(6)){
		    redraw = true;
		}
		
		else if (bounds.getWest() <= current.getWest().toFixed(6)){
		    redraw = true;
		}
		
		else if (bounds.getNorth() >= current.getNorth().toFixed(6)){
		    redraw = true;
		}
		
		else if (bounds.getEast() >= current.getEast().toFixed(6)){
		    redraw = true;
		}
		
		else {}
	    }
	    
	    if (redraw){
		var opts = { 'padding': [ 50, 50 ] };
		map.fitBounds(bounds, opts);
	    }
	}
    };
    
    return self;
})();
