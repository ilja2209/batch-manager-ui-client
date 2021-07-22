$(document).ready(function () {
    getAllProcesses();

    $('#searchButton').on('click', function (event) {
        id = $('#processIdText').val().trim()
        if (id === "") {
            getAllProcesses();
        } else {
            getProcessById(id);
        }
    });

});

function getAllProcesses() {
    $.ajax({
        url: "/api/v1/processes",
        success: function (result) {
            updateTable(result)
            //setTimeout(getAllProcesses, 7000);
        },
        error: function (data) {
            alert(data.responseText);
        }
    });
}

function getProcessById(id) {
    $.ajax({
        url: "/api/v1/processes/" + id,
        success: function (result) {
            updateTable([result])
        },
        error: function (data) {
            alert(data.responseText);
        }
    });
}

function updateTable(rows) {
    h = ""
    rows.forEach(r => {
        td = "<th scope=\"row\"><a href=\"javascript:showProcessDetails(" + r.id + ");\">" + r.id + "</th>"
        td += "<td><a href='javascript:showJsonTextArea(\"" + formatStr(r.data) + "\");'/>Request json</td>"
        td += "<td>" + r.started + "</td>"
        td += "<td>" + r.finished + "</td>"
        td += "<td>" + getBage(r.state) + "</td>"
        td += "<td><a href='javascript:showJsonTextArea(\"" + formatStr(r.context) + "\");'/>Data</td>"

        if (!isTerminatedState(r.state)) {
            td += "<td><button type=\"button\" class=\"btn btn-danger\">KILL</button></td>"
        } else {
            td += "<td></td>"
        }

        h += "<tr>" + td + "</th></tr>"
    })
    $("#main-table-body").html(h);
}

function showProcessDetails(process_id) {
    $('#processDetailsTitle').html("Process details " + process_id)
    $.ajax({
        url: "/api/v1/processes/" + process_id,
        success: function (result) {
            tasks = result.tasks == null ? [] : result.tasks;
            updateTaskTable(tasks);
            $('#processDetails').modal('show');
        }
    });
}

function updateTaskTable(tasks) {
    h = ""
    tasks.forEach(task => {
        td = "<td>" + task.id + "</td>"
        td += "<th scope=\"row\"><a target=\"_blank\" rel=\"noopener noreferrer\" href=\"http://prod-master003.offline-analytics.sel-vpc.onefactor.com:8088/proxy/" + task.application_id + "\">" + task.application_id + "</th>"
        td += "<th scope=\"row\"><a href='javascript:showJsonTextArea(\"" + formatStr(task.data) + "\");'>Data</th>"
        td += "<td>" + task.attempt + "</td>"
        td += "<td>" + task.stage_num + "</td>"
        td += "<td>" + task.started + "</td>"
        td += "<td>" + task.finished + "</td>"
        td += "<td>" + getBage(task.state) + "</td>"
        td += "<th scope=\"row\"><a href='javascript:showJsonTextArea(\"" + formatStr(task.result) + "\");'>Data</th>"
        td += "<th scope=\"row\"><a href='javascript:showJsonTextArea(\"" + formatStr(task.context) + "\");'>Data</th>"
        h += "<tr>" + td + "</th></tr>"
    })
    $("#tasks-table-body").html(h);
}

function getBage(state) {
    if (state === "SUCCEEDED" || state === "SPLIT_FINISHED") {
        return "<span class=\"badge bg-success\">" + state + "</span>";
    } else if (state === "FAILED") {
        return "<span class=\"badge bg-danger\">" + state + "</span>";
    } else if (state === "KILLING" || state === "KILLED") {
        return "<span class=\"badge bg-warning text-dark\">" + state + "</span>";
    } else {
        return "<div class=\"progress\"><div class=\"progress-bar progress-bar-striped progress-bar-animated\" role=\"progressbar\" aria-valuenow=\"100\" aria-valuemin=\"0\" aria-valuemax=\"100\" style=\"width: 100%\">" + state + "</div></div>"
    }
}

function isTerminatedState(state) {
    return state === "SUCCEEDED" || state === "KILLED" || state === "FAILED" || state === "SPLIT_FINISHED";
}

function showJsonTextArea(jsonText) {
    try {
        $('#jsonTextArea').html("");
        if (jsonText != null && jsonText != "") {
            jsonObj = JSON.parse(jsonText)
            $('#jsonTextArea').html(JSON.stringify(jsonObj, null, 4));
        }
    } catch (e) {
        $('#jsonTextArea').html(jsonText);
    }
    $('#jsonDetailsModal').modal('show');
}

function formatStr(str) {
    if (str == null) {
        return "";
    }
    return str
        .replaceAll('\\n\\t', '')
        .replaceAll('\\n', '')
        .replaceAll("'", "")
        .replaceAll('\\"', "")
        .replaceAll('"', '\\\"');
}
