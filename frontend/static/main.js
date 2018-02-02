$(document).ready(function() {
  // highlight query in the results
  // Highlighting here ensures we don't introduce unsafe characters.
  // This is not ideal and should be replaced with a template 
  // function in Go as this will break if javascript is disabled.
  $(".description").each(function(index, value){
    var q = $("#query").data("query").split(" ");
    var content = $(value).html();
    var c = content.split(" ");
    for (var i = 0; i < c.length; i++){
      for (var j = 0; j < q.length; j++){
        if (c[i].toLowerCase().indexOf(q[j].toLowerCase()) != -1){
          c[i] = "<em>" + c[i] + "</em>";
        }
      }
    }
    $(value).html(c.join(" "));
  });

  // voting
  // TODO: HMAC key
  $(document).on('click', '.arrow', function(){
    var t = $(this);
    var v = $(t).data('vote');
    var removing = false;

    if ($(t).hasClass('voted')){ // already vote for that link?
      removing = true;
      $(t).removeClass('voted');
      v = -1*v;      
    }

    d = {
      'q': $('#query').data('query'), 
      'u': $(t).parent('.vote').data('url'), 
      'v': v
    };
    
    $.ajax({
      type: "POST",
      dataType: "json",
      url: "/vote",
      data: d
    }).done(function(data) {
      $(t).siblings('.arrow').removeClass('voted'); // remove prior vote if it is different
      if (removing != true){
        $(t).addClass('voted');
      }
    }).fail(function(data) {
      $(t).siblings('.arrow').removeClass('voted');
    });
  });

  // Traditional Pagination
  $(document).on('click', '.pagination', function(){
    window.location.href = window.location.pathname + replaceQueryParam(queryString(), 'p', $(this).data('page'));
  });

  // autocomplete
  $(function(){
    $("#query").autocomplete({
      delay: 75,
      minLength: 1,
      messages: {
        noResults: "",
        results: function() {}
      },
      open: function() {
        $("ul.ui-menu").innerWidth($(this).innerWidth()); // width of input including button
      },
      source: function(request, callback){
        $.getJSON('/autocomplete', {q: request.term}, function(data){ // '{q: request.term}' changes it from ?term=b to ?q=b so nginx doesn't log query.
            callback(data.suggestions);
        });
      },
      select: function(event, ui){
        $("#query").val(ui.item.label);
        document.getElementById('form').submit();
        return false;
      },
      focus: function(event, ui) {
        $("#query").val(ui.item.label);
        return false;
      },
      }).data('ui-autocomplete')._renderItem = function(ul, item){
        var re = new RegExp(this.term, 'i');
        var re = new RegExp("^" + this.term);
        var r = item.label.replace(re, "<span style='font-weight:normal;'>" + "$&" + "</span>");
        return $("<li></li>" ).data("item.autocomplete", item).append("<a>" + r + "</a>").appendTo(ul);
      };
  });

  // fix the size of the autocomplete dropdown menu to match the size of the input
  jQuery.ui.autocomplete.prototype._resizeMenu = function(){
    var ul = this.menu.element;
    ul.outerWidth($("#query").outerWidth(true)-40); // 40 is the width of our button
  }

  // show the contributors to IA
  $("#moreinfo").on("click", function(){
    $(this).fadeOut(200, function(){
      $("#contributors").fadeIn(100);
    });
  });

  // redirect "did you mean?" queries
  $("#alternative").on("click", function(){    
    window.location.href = window.location.pathname + replaceQueryParam(queryString(), "q", $(this).attr("data-alternative"));
  });

  function queryString(){
    return window.location.search;
  }

  function replaceQueryParam(qs, param, newval) {
    var regex = new RegExp("([?;&])" + param + "[^&;]*[;&]?");
    var query = qs.replace(regex, "$1").replace(/&$/, '');
    return (query.length > 2 ? query + "&" : "?") + (newval ? param + "=" + newval : '');
  }
});
