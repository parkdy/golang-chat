$(document).ready(function() {
  var wsUrl = "ws://localhost:8080/ws";
  var ws = new WebSocket(wsUrl);

  var $chatForm = $("#chat-form"),
      $chatBox = $("#chat-box"),
      $chatLog = $("#chat-log");

  $chatForm.submit(function(event) {
    var message = $chatBox.val();
  
    $.ajax({
      url: "/messages",
      type: "POST",
      dataType: "json",
      data: { message: message }
    });

    $chatBox.val("");
    event.preventDefault();
  });

  ws.onmessage = function(event) {
    var message = event.data;
    var chatLogText = $chatLog.val();
    chatLogText = chatLogText + message + "\n";
    $chatLog.val(chatLogText);
    $chatLog.scrollTop($chatLog[0].scrollHeight);
  };
});