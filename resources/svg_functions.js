"use strict";
var details,
  searchbtn,
  unzoombtn,
  matchedtxt,
  svg,
  searching,
  currentSearchTerm,
  ignorecase,
  ignorecaseBtn;

function init(evt) {
    svg = document.getElementsByTagName("svg")[0];
    details = document.getElementById("details").firstChild;

    var blocks = document.getElementsByClassName("block");

    Array.from(blocks).forEach(function(element) {
        element.addEventListener('mouseover', block_mouseover);
        element.addEventListener('mouseout', block_mouseout);
    });
}

function block_mouseover(e) {
  var info = block_to_id(e.currentTarget);
  details.nodeValue = "Details: Block " + info;
}

function block_mouseout() {
  details.nodeValue = "Details: ";
}

// functions
function find_child(parent, name, attr) {
  var children = parent.childNodes;
  for (var i = 0; i < children.length; i++) {
    if (children[i].tagName == name)
      return attr != undefined
        ? children[i].attributes[attr].value
        : children[i];
  }
  return;
}

function find_group(node) {
  var parent = node.parentElement;
  if (!parent) return;
  if (parent.id == "frames") return node;
  return find_group(parent);
}

function block_to_id(node) {
  return node.id.split("_").at(-1)
}