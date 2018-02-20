import { Component, OnInit } from "@angular/core";
import { AutomatedUIServices } from "../automated-ui-services";
import { ToastsManager } from "ng2-toastr/ng2-toastr";
@Component({
  selector: "app-test-suites",
  templateUrl: "./test-suites.component.html",
  styleUrls: ["./test-suites.component.css"]
})
export class TestSuitesComponent {
  testSuiteData = {};
  testCaseData = {};
  selectedTestCaseData = {};
  testSuitePath = undefined;
  testSuiteFileName = undefined;
  testSuiteSchema = { layout: true };
  constructor(
    private automatedUIServices: AutomatedUIServices,
    private toastr: ToastsManager
  ) { }

  ngOnInit() {
    this.automatedUIServices
      .getSchema$("assets/testSuite_schema.json")
      .subscribe((data: any) => {
        this.testSuiteSchema = data;
      });
  }

  onAdd() {
    this.automatedUIServices.getAllTestSuite$(this.testSuitePath).subscribe(
      (data: any) => {
        this.testSuiteData = data;

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
              "No Test Suite Directory added",
              "An error occurred!"
            );
            break;
          }
          default: {
            this.toastr.error("An error occurred!");
          }
        }
      }
    );
  }

  getTestCases() {
    this.automatedUIServices
      .getAllTestCases$(this.testSuitePath)
      .subscribe((data: any) => {
        this.testCaseData = data;
      });
  }


  onDelete() {

  }

  onCancel() {
    //clear schema and moving info back into available (get method)
    this.automatedUIServices
      .getTestSuite$(this.testSuitePath, this.testSuiteFileName)
      .subscribe(
      data => {
        this.testSuiteData = data;
        this.toastr.success("Previous data reloaded!");
      },
      error => {
        this.testSuiteData = undefined;
      }
      );
  }

  onSave(testSuiteData) {
    this.automatedUIServices
      .postTestSuite$(testSuiteData, this.testSuitePath)
      .subscribe(
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
            this.toastr.error(
              "Your data did not save!",
              "An error occurred!"
            );
          }
        }
      }
      );
  }

  onSelectAll() {
    this.automatedUIServices
      .getAllTestCases$(this.testSuitePath)
      .subscribe((data: any) => {
        this.selectedTestCaseData = data;
      });
  }
  onSelectOne() { }
  onReverseOne() { }
  onReverseAll() {
    this.selectedTestCaseData = undefined;
  }
}
