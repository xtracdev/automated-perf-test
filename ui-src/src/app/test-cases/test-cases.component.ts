import { Component, OnInit } from "@angular/core";
import { AutomatedUIServices } from "../automated-ui-services";
import { ToastsManager } from "ng2-toastr/ng2-toastr";
import { JsonSchemaFormModule } from "angular2-json-schema-form";
import { HttpClient } from "@angular/common/http";
import "rxjs/add/operator/map";

@Component({
  selector: "app-test-cases",
  templateUrl: "./test-cases.component.html",
  styleUrls: ["./test-cases.component.css"]
})
export class TestCasesComponent implements OnInit {
  testCaseData = {};
  testCasePath = undefined;
  testCaseSchema = { layout: true };

  constructor(
    private automatedUIServices: AutomatedUIServices,
    private toastr: ToastsManager,
    private http: HttpClient
  ) { }

  ngOnInit() {
    this.automatedUIServices
      .getSchema$("assets/testCase_schema.json")
      .subscribe((data: any) => {
        this.testCaseSchema = data;
      });
  }


  onAdd() {
    this.automatedUIServices.getAllCases$(this.testCasePath).subscribe(
      (data: any) => {
        this.testCaseData = data;

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

  onDelete() { }

  onSave(testCaseData) {
    this.automatedUIServices.postTestCases$(testCaseData, this.testCasePath).subscribe(
      data => {
        this.toastr.success("Your data has been saved!", "Success!");
      },

      error => {
        switch (error.status) {
          case 500: {
            this.toastr.error("An error has occurred!", "Check the logs!");
            break;
          }
          case 409: {
            this.toastr.error("File already exists!", "An error occurred!");
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
  onCancel() { }
}
