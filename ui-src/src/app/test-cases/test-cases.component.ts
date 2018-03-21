import { Component, OnInit } from "@angular/core";
import { ConfigurationService } from "../configurations/configuration.service";
import { TestCaseService } from "./test-case.service";
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
  testCaseFileName = undefined;
  testCases = [];


  constructor(
    private configurationService: ConfigurationService,
    private testCaseService: TestCaseService,
    private toastr: ToastsManager,
    private http: HttpClient
  ) { }

  ngOnInit() {
    this.configurationService
      .getSchema$("assets/testCase_schema.json")
      .subscribe((data: any) => {
        this.testCaseSchema = data;
      });
  }

  onLoad() {
    this.testCaseService.getTestCases$(this.testCasePath).subscribe(
      (data: any) => {
        this.testCases = data;
        this.toastr.success("Your Test Cases have loaded!", "Success!");
      },
      error => {
        switch (error.status) {
          case 500: {
            this.toastr.error("An error has occurred!", "Check the logs!");
            break;
          }
          case 400: {
            this.toastr.error(
              "No Test Case Directory added",
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

  onSelectCase(testCase, i) {
    this.testCaseData = testCase;
    this.testCaseFileName = testCase.testname;
  }

  onAdd() {
    this.testCaseData = undefined;
    this.testCaseFileName = undefined;
  }

  onDelete() { }

  onSave(testCaseData) {
    this.testCaseService.postTestCase$(testCaseData, this.testCasePath).subscribe(
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


  onUpdate(testCaseData) {
    this.testCaseService
      .putTestCase$(testCaseData, this.testCasePath, this.testCaseFileName)
      .subscribe(
      data => {
        this.toastr.success("Success!");
        this.onLoad();
      },
      error => {
        switch (error.status) {
          case 404: {
            this.toastr.error("File not found", "An error occured!");

            break;
          }
          case 400: {
            this.toastr.error(
              "File must be specified!",
              "An error occurred!"
            );
            break;
          }
          case 500: {
            this.toastr.error("Internal server error!");
            break;
          }
          default: {
            this.toastr.error("File was not updated!", "An error occurred!");
          }
        }
      }
      );
  }

  onCancel() {
    this.testCaseService
      .getTestCase$(this.testCasePath, this.testCaseFileName)
      .subscribe(
      data => {
        this.testCaseData = data;
        this.toastr.success("Previous data reloaded!");
      },
      error => {
        this.onAdd();
      }
      );
  }
}
