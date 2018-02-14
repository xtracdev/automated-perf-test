import { Component, OnInit } from "@angular/core";
import { AutomatedUIServices } from "../automated-ui-services";
import { ToastsManager } from "ng2-toastr/ng2-toastr";
@Component({
  selector: "app-test-suites",
  templateUrl: "./test-suites.component.html",
  styleUrls: ["./test-suites.component.css"]
})
export class TestSuitesComponent {

  testSuitePath = undefined;
  constructor(
    private automatedUIServices: AutomatedUIServices,
    private toastr: ToastsManager
  ) { }
  testSchema = {
    type: "object",
    properties: {
      name: { type: "string" },
      testStrategy: {
        type: "string",
        "enum": [
          "SuiteBased",
          "ServiceBased"
        ]
      },
      description: { type: "string" }

    }
  }

  onSubmit(testSuiteData) {
    this.automatedUIServices.postTestSuite$(testSuiteData, this.testSuitePath).subscribe(
      data => {
        this.toastr.success("Your data has been saved!", "Success!");
      },

      error => {
        switch (error.status) {
          case 500: {
            this.toastr.error("An error has occurred!", "Check the logs!");
            break;
          }
          case 400: {
            this.toastr.error(
              "Some of the fields do not conform to the schema!",
              "An error occurred!"
            );
            break;
          }
          default: {
            this.toastr.error("Your data did not save!", "An error occurred!");
          }
        }
      }
    );
  }


  onDelete() { }
  onCancel() { }
  onSave() { }
  onAdd(){}
}
