/*
   Licensed to the Apache Software Foundation (ASF) under one or more
   contributor license agreements.  See the NOTICE file distributed with
   this work for additional information regarding copyright ownership.
   The ASF licenses this file to You under the Apache License, Version 2.0
   (the "License"); you may not use this file except in compliance with
   the License.  You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
var showControllersOnly = false;
var seriesFilter = "";
var filtersOnlySampleSeries = true;

/*
 * Add header in statistics table to group metrics by category
 * format
 *
 */
function summaryTableHeader(header) {
    var newRow = header.insertRow(-1);
    newRow.className = "tablesorter-no-sort";
    var cell = document.createElement('th');
    cell.setAttribute("data-sorter", false);
    cell.colSpan = 1;
    cell.innerHTML = "Requests";
    newRow.appendChild(cell);

    cell = document.createElement('th');
    cell.setAttribute("data-sorter", false);
    cell.colSpan = 3;
    cell.innerHTML = "Executions";
    newRow.appendChild(cell);

    cell = document.createElement('th');
    cell.setAttribute("data-sorter", false);
    cell.colSpan = 7;
    cell.innerHTML = "Response Times (ms)";
    newRow.appendChild(cell);

    cell = document.createElement('th');
    cell.setAttribute("data-sorter", false);
    cell.colSpan = 1;
    cell.innerHTML = "Throughput";
    newRow.appendChild(cell);

    cell = document.createElement('th');
    cell.setAttribute("data-sorter", false);
    cell.colSpan = 2;
    cell.innerHTML = "Network (KB/sec)";
    newRow.appendChild(cell);
}

/*
 * Populates the table identified by id parameter with the specified data and
 * format
 *
 */
function createTable(table, info, formatter, defaultSorts, seriesIndex, headerCreator) {
    var tableRef = table[0];

    // Create header and populate it with data.titles array
    var header = tableRef.createTHead();

    // Call callback is available
    if(headerCreator) {
        headerCreator(header);
    }

    var newRow = header.insertRow(-1);
    for (var index = 0; index < info.titles.length; index++) {
        var cell = document.createElement('th');
        cell.innerHTML = info.titles[index];
        newRow.appendChild(cell);
    }

    var tBody;

    // Create overall body if defined
    if(info.overall){
        tBody = document.createElement('tbody');
        tBody.className = "tablesorter-no-sort";
        tableRef.appendChild(tBody);
        var newRow = tBody.insertRow(-1);
        var data = info.overall.data;
        for(var index=0;index < data.length; index++){
            var cell = newRow.insertCell(-1);
            cell.innerHTML = formatter ? formatter(index, data[index]): data[index];
        }
    }

    // Create regular body
    tBody = document.createElement('tbody');
    tableRef.appendChild(tBody);

    var regexp;
    if(seriesFilter) {
        regexp = new RegExp(seriesFilter, 'i');
    }
    // Populate body with data.items array
    for(var index=0; index < info.items.length; index++){
        var item = info.items[index];
        if((!regexp || filtersOnlySampleSeries && !info.supportsControllersDiscrimination || regexp.test(item.data[seriesIndex]))
                &&
                (!showControllersOnly || !info.supportsControllersDiscrimination || item.isController)){
            if(item.data.length > 0) {
                var newRow = tBody.insertRow(-1);
                for(var col=0; col < item.data.length; col++){
                    var cell = newRow.insertCell(-1);
                    cell.innerHTML = formatter ? formatter(col, item.data[col]) : item.data[col];
                }
            }
        }
    }

    // Add support of columns sort
    table.tablesorter({sortList : defaultSorts});
}

$(document).ready(function() {

    // Customize table sorter default options
    $.extend( $.tablesorter.defaults, {
        theme: 'blue',
        cssInfoBlock: "tablesorter-no-sort",
        widthFixed: true,
        widgets: ['zebra']
    });

    var data = {"OkPercent": 68.9375, "KoPercent": 31.0625};
    var dataset = [
        {
            "label" : "FAIL",
            "data" : data.KoPercent,
            "color" : "#FF6347"
        },
        {
            "label" : "PASS",
            "data" : data.OkPercent,
            "color" : "#9ACD32"
        }];
    $.plot($("#flot-requests-summary"), dataset, {
        series : {
            pie : {
                show : true,
                radius : 1,
                label : {
                    show : true,
                    radius : 3 / 4,
                    formatter : function(label, series) {
                        return '<div style="font-size:8pt;text-align:center;padding:2px;color:white;">'
                            + label
                            + '<br/>'
                            + Math.round10(series.percent, -2)
                            + '%</div>';
                    },
                    background : {
                        opacity : 0.5,
                        color : '#000'
                    }
                }
            }
        },
        legend : {
            show : true
        }
    });

    // Creates APDEX table
    createTable($("#apdexTable"), {"supportsControllersDiscrimination": true, "overall": {"data": [0.64565625, 500, 1500, "Total"], "isController": false}, "titles": ["Apdex", "T (Toleration threshold)", "F (Frustration threshold)", "Label"], "items": [{"data": [0.832, 500, 1500, "Get All Achievements"], "isController": false}, {"data": [0.842, 500, 1500, "Get User Profile"], "isController": false}, {"data": [0.646, 500, 1500, "Get All Unlocked Achievements"], "isController": false}, {"data": [0.846, 500, 1500, "Update Profile"], "isController": false}, {"data": [0.647, 500, 1500, "Comment on Post"], "isController": false}, {"data": [0.2095, 500, 1500, "Get Schedules"], "isController": false}, {"data": [0.801, 500, 1500, "Login"], "isController": false}, {"data": [0.686, 500, 1500, "Get Exercises"], "isController": false}, {"data": [0.841, 500, 1500, "Get Unlocked Achievements"], "isController": false}, {"data": [0.672, 500, 1500, "Schedule Routine"], "isController": false}, {"data": [0.582, 500, 1500, "Create Post"], "isController": false}, {"data": [0.848, 500, 1500, "Get Following"], "isController": false}, {"data": [0.146, 500, 1500, "Get Posts"], "isController": false}, {"data": [0.217, 500, 1500, "Get Routines"], "isController": false}, {"data": [0.851, 500, 1500, "Get Followers"], "isController": false}, {"data": [0.664, 500, 1500, "Create Routine"], "isController": false}]}, function(index, item){
        switch(index){
            case 0:
                item = item.toFixed(3);
                break;
            case 1:
            case 2:
                item = formatDuration(item);
                break;
        }
        return item;
    }, [[0, 0]], 3);

    // Create statistics table
    createTable($("#statisticsTable"), {"supportsControllersDiscrimination": true, "overall": {"data": ["Total", 16000, 4970, 31.0625, 259.6616249999995, 0, 11605, 17.0, 167.6999999999989, 1739.949999999999, 5267.939999999999, 154.2510628862302, 6590.862029202859, 78.45988222087787], "isController": false}, "titles": ["Label", "#Samples", "FAIL", "Error %", "Average", "Min", "Max", "Median", "90th pct", "95th pct", "99th pct", "Transactions/s", "Received", "Sent"], "items": [{"data": ["Get All Achievements", 1000, 168, 16.8, 14.885000000000002, 1, 90, 9.0, 38.0, 50.0, 69.0, 10.087661780875811, 8.068612334940633, 2.954581563789329], "isController": false}, {"data": ["Get User Profile", 1000, 158, 15.8, 17.727999999999984, 1, 105, 12.0, 42.0, 52.0, 71.99000000000001, 10.07688664510213, 28.84141743691869, 2.8640007633241633], "isController": false}, {"data": ["Get All Unlocked Achievements", 1000, 354, 35.4, 14.363000000000003, 1, 111, 9.0, 39.89999999999998, 53.0, 73.98000000000002, 10.079121100640023, 3.6154516076198155, 3.0071313718943706], "isController": false}, {"data": ["Update Profile", 1000, 154, 15.4, 19.756999999999987, 1, 101, 13.0, 47.0, 57.0, 82.0, 10.079121100640023, 28.993989536612407, 31.257165625157484], "isController": false}, {"data": ["Comment on Post", 1000, 353, 35.3, 16.994000000000003, 1, 97, 12.0, 44.0, 55.94999999999993, 80.95000000000005, 10.078308456708626, 3.7686377843342775, 3.9286506152807315], "isController": false}, {"data": ["Get Schedules", 1000, 591, 59.1, 847.7869999999999, 1, 6098, 104.0, 2802.8999999999996, 3786.6499999999996, 5203.950000000001, 9.659036028204387, 110.89582654182362, 2.7908199917898195], "isController": false}, {"data": ["Login", 1000, 199, 19.9, 76.92000000000004, 0, 178, 84.0, 117.0, 127.0, 148.0, 10.016828271496115, 29.420823467800908, 3.9015350476300186], "isController": false}, {"data": ["Get Exercises", 1000, 314, 31.4, 27.790999999999997, 1, 108, 22.0, 59.0, 70.0, 89.99000000000001, 9.734349599431514, 6062.064349830623, 2.8125806125826203], "isController": false}, {"data": ["Get Unlocked Achievements", 1000, 159, 15.9, 16.58599999999997, 1, 93, 11.0, 40.0, 52.0, 72.99000000000001, 10.087051252307413, 3.1637583123606725, 3.0359660116808054], "isController": false}, {"data": ["Schedule Routine", 1000, 328, 32.8, 33.761000000000024, 1, 143, 23.0, 76.0, 90.94999999999993, 113.98000000000002, 9.659689151203114, 4.350671167276837, 4.627924999154778], "isController": false}, {"data": ["Create Post", 1000, 418, 41.8, 24.610999999999983, 1, 144, 15.0, 62.0, 73.94999999999993, 105.99000000000001, 10.078105316200554, 4.799638605442177, 4.178673784328547], "isController": false}, {"data": ["Get Following", 1000, 152, 15.2, 14.859000000000005, 1, 90, 9.0, 36.89999999999998, 47.94999999999993, 65.98000000000002, 10.078511605406113, 3.1433735990868867, 3.0004339271676357], "isController": false}, {"data": ["Get Posts", 1000, 578, 57.8, 2040.7090000000028, 1, 11605, 753.5, 6144.2, 8756.099999999997, 10639.6, 10.074145712443585, 245.97929432498188, 2.89354256704344], "isController": false}, {"data": ["Get Routines", 1000, 559, 55.9, 938.2629999999996, 1, 6319, 139.5, 3097.2, 4376.049999999997, 5695.680000000002, 9.733781087263347, 121.50158021852339, 2.8733665498612937], "isController": false}, {"data": ["Get Followers", 1000, 149, 14.9, 15.846000000000007, 1, 107, 10.0, 40.0, 52.0, 73.0, 10.079222690346121, 3.137541939015663, 3.000645621459673], "isController": false}, {"data": ["Create Routine", 1000, 336, 33.6, 33.72599999999995, 1, 127, 24.0, 75.0, 88.94999999999993, 114.99000000000001, 9.736339915099116, 4.364523546121042, 4.206669331989718], "isController": false}]}, function(index, item){
        switch(index){
            // Errors pct
            case 3:
                item = item.toFixed(2) + '%';
                break;
            // Mean
            case 4:
            // Mean
            case 7:
            // Median
            case 8:
            // Percentile 1
            case 9:
            // Percentile 2
            case 10:
            // Percentile 3
            case 11:
            // Throughput
            case 12:
            // Kbytes/s
            case 13:
            // Sent Kbytes/s
                item = item.toFixed(2);
                break;
        }
        return item;
    }, [[0, 0]], 0, summaryTableHeader);

    // Create error table
    createTable($("#errorsTable"), {"supportsControllersDiscrimination": false, "titles": ["Type of error", "Number of errors", "% in errors", "% in all samples"], "items": [{"data": ["503/Service Unavailable", 1744, 35.09054325955734, 10.9], "isController": false}, {"data": ["Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 1284, 25.835010060362173, 8.025], "isController": false}, {"data": ["500/Internal Server Error", 1824, 36.70020120724346, 11.4], "isController": false}, {"data": ["401/Unauthorized", 118, 2.374245472837022, 0.7375], "isController": false}]}, function(index, item){
        switch(index){
            case 2:
            case 3:
                item = item.toFixed(2) + '%';
                break;
        }
        return item;
    }, [[1, 1]]);

        // Create top5 errors by sampler
    createTable($("#top5ErrorsBySamplerTable"), {"supportsControllersDiscrimination": false, "overall": {"data": ["Total", 16000, 4970, "500/Internal Server Error", 1824, "503/Service Unavailable", 1744, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 1284, "401/Unauthorized", 118, "", ""], "isController": false}, "titles": ["Sample", "#Samples", "#Errors", "Error", "#Errors", "Error", "#Errors", "Error", "#Errors", "Error", "#Errors", "Error", "#Errors"], "items": [{"data": ["Get All Achievements", 1000, 168, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 80, "503/Service Unavailable", 75, "500/Internal Server Error", 13, "", "", "", ""], "isController": false}, {"data": ["Get User Profile", 1000, 158, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 79, "503/Service Unavailable", 64, "500/Internal Server Error", 15, "", "", "", ""], "isController": false}, {"data": ["Get All Unlocked Achievements", 1000, 354, "500/Internal Server Error", 209, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 77, "503/Service Unavailable", 68, "", "", "", ""], "isController": false}, {"data": ["Update Profile", 1000, 154, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 79, "503/Service Unavailable", 61, "500/Internal Server Error", 14, "", "", "", ""], "isController": false}, {"data": ["Comment on Post", 1000, 353, "500/Internal Server Error", 212, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 76, "503/Service Unavailable", 65, "", "", "", ""], "isController": false}, {"data": ["Get Schedules", 1000, 591, "503/Service Unavailable", 311, "500/Internal Server Error", 196, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 84, "", "", "", ""], "isController": false}, {"data": ["Login", 1000, 199, "401/Unauthorized", 118, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 81, "", "", "", "", "", ""], "isController": false}, {"data": ["Get Exercises", 1000, 314, "500/Internal Server Error", 156, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 84, "503/Service Unavailable", 74, "", "", "", ""], "isController": false}, {"data": ["Get Unlocked Achievements", 1000, 159, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 80, "503/Service Unavailable", 67, "500/Internal Server Error", 12, "", "", "", ""], "isController": false}, {"data": ["Schedule Routine", 1000, 328, "500/Internal Server Error", 172, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 83, "503/Service Unavailable", 73, "", "", "", ""], "isController": false}, {"data": ["Create Post", 1000, 418, "500/Internal Server Error", 222, "503/Service Unavailable", 119, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 77, "", "", "", ""], "isController": false}, {"data": ["Get Following", 1000, 152, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 79, "503/Service Unavailable", 56, "500/Internal Server Error", 17, "", "", "", ""], "isController": false}, {"data": ["Get Posts", 1000, 578, "503/Service Unavailable", 264, "500/Internal Server Error", 236, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 78, "", "", "", ""], "isController": false}, {"data": ["Get Routines", 1000, 559, "503/Service Unavailable", 295, "500/Internal Server Error", 180, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 84, "", "", "", ""], "isController": false}, {"data": ["Get Followers", 1000, 149, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 79, "503/Service Unavailable", 53, "500/Internal Server Error", 17, "", "", "", ""], "isController": false}, {"data": ["Create Routine", 1000, 336, "500/Internal Server Error", 153, "503/Service Unavailable", 99, "Non HTTP response code: java.net.BindException/Non HTTP response message: Address already in use: getsockopt", 84, "", "", "", ""], "isController": false}]}, function(index, item){
        return item;
    }, [[0, 0]], 0);

});
