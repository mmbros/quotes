class Item {
    constructor(d) {
        var className = function (d) {
            switch (d.status) {
                case 2: return "success";
                case 3: return "error";
                case 4: return "canceled";
            }
            return d.status;
        };
        var result = function (d) {
            if (d.status == 2) {
                return `${d.price} ${d.currency}`;
            }
            return "n/a";
        };
        this.task_id = d.isin;
        this.status = className(d);
        this.worker_id = d.source;
        this.worker_inst = d.instance;
        this.msec_start = (new Date(d.time_start)).getTime();
        this.msec_end = (new Date(d.time_end)).getTime();
        this.result = result(d);
        this.error = d.error;
    }

    worker() {
        return `${this.worker_id}[${this.worker_inst}]`
    }

    tooltip() {
        var s = "task: " + this.task_id +
            "\nworker: " + this.worker() +
            "\nstart: " + this.msec_start +
            "\nend: " + this.msec_end;
        switch (this.status) {
            case "success":
                s += "\nresult: " + this.result; break;
            case "error":
                s += "\nerror: " + this.error; break;
        }
        return s;
    }
}



function onlyUnique(value, index, self) {
    return self.indexOf(value) === index;
};



function json2Scenario(jsonData) {

    var items = jsonData.map(d => new Item(d));

    // get min and max timestamp
    var msec_min = d3.min(items, d => d.msec_start);
    var msec_max = d3.max(items, d => d.msec_end);

    // shift timestamps to msec_min
    items = items.map(function (d) {
        d.msec_start -= msec_min;
        d.msec_end -= msec_min;
        return d;
    })

    // tasks results
    var tasks = {}
    items.forEach(item => {
        var task_id = item.task_id;
        var status = item.status;

        if (tasks[task_id] === undefined) {
            tasks[task_id] = status
        } else {
            switch (status) {
                case "success":
                    tasks[task_id] = status;
                case "error":
                    if (tasks[task_id] != "success") {
                        tasks[task_id] = status;
                    }
            }
        }
    });

    return {
        msec_min: msec_min,
        msec_max: msec_max,
        workers: d3.sort(d3.map(items, d => d.worker()).filter(onlyUnique)),
        tasks: tasks,
        items: items,
    }
};


function showInfo(scenario) {
    showInfoJob(); // clear info

    var tasks_success = 0;
    Object.values(scenario.tasks).forEach(function (status) {
        if (status == "success") tasks_success++;
    })

    var tasks_error = 0;
    Object.values(scenario.tasks).forEach(function (status) {
        if (status == "error") tasks_error++;
    })

    var msecs = {}
    scenario.items.forEach(function (item) {

        if (item.status == "canceled")
            return;

        var msec = msecs[item.task_id];
        if (msec === undefined) {
            msecs[item.task_id] = {
                min: item.msec_start,
                max: item.msec_end,
            }
        } else {
            if (msec.min > item.msec_start) {
                msec.min = item.msec_start
            }
            if (msec.max < item.msec_end) {
                msec.max = item.msec_end
            }
            msecs[item.task_id] = msec
        }
    })

    var elapsed_min = 9999999999;
    var elapsed_max = 0;
    var elapsed_sum = 0;

    Object.values(msecs).forEach(function (msec) {
        var elapsed = msec.max - msec.min;
        if (elapsed_min > elapsed) {
            elapsed_min = elapsed
        }
        if (elapsed_max < elapsed) {
            elapsed_max = elapsed
        }
        elapsed_sum += elapsed;
    })

    var tasks_tot = Object.keys(scenario.tasks).length;

    var array = [
        ["workers", `${scenario.workers.length}`],
        ["tasks", `${tasks_tot}`],
        ["tasks_success", `${tasks_success}`],
        ["tasks_error", `${tasks_error}`],
        ["elapsed_tot", `${scenario.msec_max - scenario.msec_min} ms`],
        ["elapsed_task_min", `${elapsed_min} ms`],
        ["elapsed_task_avg", `${Math.round(elapsed_sum / tasks_tot)} ms`],
        ["elapsed_task_max", `${elapsed_max} ms`]
    ];
    array.forEach(function (d) {
        document.getElementById(d[0]).innerHTML = d[1];
    });
};

function showInfoJob(data) {
    if (data === undefined) {
        // clear
        const array = ["task_id", "task_status",
            "task_workers", "task_workers_success", "task_workers_error", "task_workers_canceled",
            "task_msec_min", "task_msec_max", "task_msec_elapsed",
            "job_task_id", "job_worker_id", "job_status", "job_msec_start", "job_msec_end", "job_msec_elapsed"]
        array.forEach(function (key) {
            document.getElementById(key).innerHTML = "";
        });
        return;
    }

    var task_workers = 0;
    var msec_min = 999999999;
    var msec_max = 0;
    var status = "canceled";
    var task_status = {};

    task_status["success"] = 0;
    task_status["error"] = 0;
    task_status["canceled"] = 0;

    d3.select("#graphTasks")
        .selectAll("rect")
        .nodes()
        .map(function (d) {
            var job = d3.select(d).datum();
            if (job.task_id == data.task_id) {
                task_workers++;
                task_status[job.status]++;

                switch (job.status) {
                    case "success":
                        status = job.status;
                    case "error":
                        if (status != "success") {
                            status = job.status;
                        }
                }

                if (job.status != "canceled") {
                    if (msec_min > job.msec_start) {
                        msec_min = job.msec_start
                    }
                    if (msec_max < job.msec_end) {
                        msec_max = job.msec_end
                    }
                }
            }
        })

    var array = [
        ["task_id", data.task_id],
        ["task_status", status],
        ["task_workers", task_workers],
        ["task_workers_success", task_status["success"]],
        ["task_workers_error", task_status["error"]],
        ["task_workers_canceled", task_status["canceled"]],
        ["task_msec_min", `${msec_min} ms`],
        ["task_msec_max", `${msec_max} ms`],
        ["task_msec_elapsed", `${msec_max - msec_min} ms`],
        ["job_task_id", data.task_id],
        ["job_worker_id", data.worker_id],
        ["job_status", data.status],
        ["job_msec_start", `${data.msec_start} ms`],
        ["job_msec_end", `${data.msec_end} ms`],
        ["job_msec_elapsed", `${data.msec_end - data.msec_start} ms`]
    ];

    array.forEach(function (d) {
        document.getElementById(d[0]).innerHTML = d[1];
    });
}


function drawGraphWorkers(scenario) {
    var margin = { top: 10, right: 40, bottom: 30, left: 100 },
        width = 1600 - margin.left - margin.right,
        height = 800 - margin.top - margin.bottom;

    // remove old graph
    d3.select("#graphWorkers").select("svg").remove();

    // append the svg object to the body of the page
    var svG = d3.select("#graphWorkers")
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform",
            "translate(" + margin.left + "," + margin.top + ")");

    // X scale and Axis
    var x = d3.scaleLinear()
        .domain([0, scenario.msec_max - scenario.msec_min])
        .range([0, width]);
    svG
        .append('g')
        .attr("transform", "translate(0," + height + ")")
        .call(d3.axisBottom(x));

    // Y scale and Axis
    var y = d3.scaleBand()
        .domain(scenario.workers)
        .range([height, 0])
        .padding([0.2]);
    svG
        .append('g')
        .call(d3.axisLeft(y));

    // draw items
    svG
        .selectAll("whatever")
        .data(scenario.items)
        .enter()
        .append("rect")
        .attr("x", d => x(d.msec_start))
        .attr("y", d => y(d.worker()))
        .attr("width", d => x(d.msec_end - d.msec_start))
        .attr("height", y.bandwidth())
        .attr("class", d => d.status)
        .on("click", graph_onclick)
        .append("title")
        .text(d => d.tooltip())
        ;

};

function drawGraphTasks(scenario) {
    var margin = { top: 10, right: 40, bottom: 30, left: 100 },
        width = 1600 - margin.left - margin.right,
        height = 800 - margin.top - margin.bottom;

    // remove old graph
    d3.select("#graphTasks").select("svg").remove();

    // append the svg object to the body of the page
    var svG = d3.select("#graphTasks")
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform",
            "translate(" + margin.left + "," + margin.top + ")");

    var taskids = d3.sort(Object.keys(scenario.tasks));

    // X scale and Axis
    var x = d3.scaleLinear()
        .domain([0, scenario.msec_max - scenario.msec_min])
        .range([0, width]);
    svG
        .append('g')
        .attr("transform", "translate(0," + height + ")")
        .call(d3.axisBottom(x));

    // Y scale and Axis
    var y = d3.scaleBand()
        .domain(taskids)
        .range([height, 0])
    // .padding([0.2]);
    svG
        .append('g')
        .call(d3.axisLeft(y)
            // .tickValues(["t010", "t020"])
        );

    // draw items
    svG
        .selectAll("whatever")
        .data(scenario.items)
        .enter()
        .append("rect")
        .attr("x", d => x(d.msec_start))
        .attr("y", d => y(d.task_id))
        .attr("width", d => x(d.msec_end - d.msec_start))
        .attr("height", y.bandwidth())
        .attr("class", d => d.status)
        .on("click", graph_onclick)
        .append("title")
        .text(d => d.tooltip())
        ;
};

function jsonGet(filename) {
    d3.json("/data/" + filename).then(
        function (jsonData) {
            var scenario = json2Scenario(jsonData);
            drawGraphWorkers(scenario);
            drawGraphTasks(scenario);
            showInfo(scenario);
        }
    );
};

function select_onchange() {
    jsonGet(this.value);
};

function graph_onclick(event, data) {
    var fn = function (d) {
        var task = d3.select(d);
        var datum = task.datum();
        task.classed("highlight", datum.task_id == data.task_id);
        task.classed("current", (datum.task_id == data.task_id) && (datum.worker_id == data.worker_id));
    }

    d3.select("#graphWorkers").selectAll("rect").nodes().map(fn);
    d3.select("#graphTasks").selectAll("rect").nodes().map(fn);


    showInfoJob(data);
};

function initSelect(selector) {
    var theSelect = d3.select(selector)

    // bind select change event
    theSelect.on('change', select_onchange);

    // fill the select with the data json files
    d3.json("/data").then(
        function (data) {
            // console.log(data);
            theSelect
                .selectAll('option')
                .data(data)
                .join("option")
                .text(d => d)
                .attr("value", d => d)
                ;
            jsonGet(theSelect.property('value'));
        }
    )
};

initSelect("#filejson");
