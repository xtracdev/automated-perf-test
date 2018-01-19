import { Component, OnInit } from "@angular/core";
import { Http } from "@angular/http";
import { AutomatedUIServices } from "../automated-ui-services";
import { JsonSchemaFormModule } from "angular2-json-schema-form";
import { ToastsManager } from "ng2-toastr/ng2-toastr";
import { ToastOptions } from "ng2-toastr/src/toast-options";
@Component({
  selector: "app-configurations",
  templateUrl: "./configurations.component.html",
  styleUrls: ["./configurations.component.css"]
})
export class ConfigurationsComponent implements OnInit {
  formData = {};
  constructor(
    private automatedUIServices: AutomatedUIServices,
    private toastr: ToastsManager
  ) {}

  exampleSchema = {
    type: "object",
    properties: {
      apiName: { type: "string" },
      targetHost: { type: "string" },
      targetPort: {
        type: "string",
        minimum: 1,
        maximum: 65535
      },
      memoryEndpoint: { type: "string" },
      numIterations: { type: "integer", minimum: 1 },
      concurrentUsers: { type: "integer", minimum: 1 },
      allowablePeakMemoryVariance: {
        type: "number",
        minimum: 0,
        maximum: 100
      },
      allowableServiceResponseTimeVariance: {
        type: "number",
        minimum: 0,
        maximum: 100
      },
      testSuite: {
        type: "string",
        enum: ["Default-1", "Default-2", "Default-3"]
      },
      requestDelay: { type: "integer", minimum: 1 },
      TPSFreq: { type: "integer", minimum: 1 },
      rampUsers: { type: "integer", minimum: 1 },
      rampDelay: { type: "integer", minimum: 0 },
      testCaseDir: { type: "string" },
      testSuiteDir: { type: "string" },
      baseStatsOutputDir: { type: "string" },
      reportOutputDir: { type: "string" }
    },
    required: [
      "apiName",
      "targetHost",
      "targetPort",
      "numIterations",
      "concurrentUsers",
      "allowablePeakMemoryVariance",
      "allowableServiceResponseTimeVariance",
      "testSuite",
      "requestDelay",
      "TPSFreq",
      "rampUsers",
      "rampDelay",
      "testCaseDir",
      "testSuiteDir",
      "baseStatsOutputDir",
      "reportOutputDir"
    ]
  };

  exampleData = {
    allowablePeakMemoryVariance: 15,
    allowableServiceResponseTimeVariance: 15
  };

  layout = [
    {
      type: "flex",
      "flex-flow": "row wrap",
      items: [
        "apiName",
        "numIterations",
        {
          key: "requestDelay",
          title: "Request Delay (ms)"
        }
      ]
    },
    {
      type: "flex",
      "flex-flow": "row wrap",
      items: [
        "targetHost",
        "concurrentUsers",
        {
          key: "TPSFreq",
          title: "TPS Frequency (s)"
        }
      ]
    },
    {
      type: "flex",
      "flex-flow": "row wrap",
      items: [
        "targetPort",
        {
          key: "allowablePeakMemoryVariance",
          title: "Memory Variance (%)"
        },
        "rampUsers"
      ]
    },
    {
      type: "flex",
      "flex-flow": "row wrap",
      items: [
        "memoryEndpoint",
        {
          key: "allowableServiceResponseTimeVariance",
          title: "Service Variance (%)"
        },
        {
          key: "rampDelay",
          title: "Ramp Delay (s)"
        }
      ]
    },
    { key: "testSuite" },
    { key: "testCaseDir", title: "Test Case Directory" },
    { key: "testSuiteDir", title: "Test Suites Directory" },
    { key: "baseStatsOutputDir", title: "Base Stats Output Directory" },
    { key: "reportOutputDir", title: "Report Output Directory" }
  ];

  ngOnInit() {
    this.formData = {
      allowablePeakMemoryVariance: 15,
      allowableServiceResponseTimeVariance: 15
    };
  }

  onSubmit(configData) {
    this.toastr.success("Your data has been save!", "Success!");
    this.automatedUIServices
      .createJsonFile(configData)
      .subscribe(data => console.log(data));
  }
  onCancel() {
    this.formData = {};
  }
}
