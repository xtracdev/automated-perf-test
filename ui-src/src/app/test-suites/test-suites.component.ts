import {Component, OnInit, Input, Output, EventEmitter} from "@angular/core";
import {TestSuiteService} from "./test-suite.service";
import {TestCaseService} from "../test-cases/test-case.service";
import {ConfigurationService} from "../configurations/configuration.service";
import {ToastsManager} from "ng2-toastr/ng2-toastr";
import {TestCasesSelectionComponent} from "../shared/test-cases-selection/test-cases-selection.component";
@Component({
  selector: "app-test-suites",
  templateUrl: "./test-suites.component.html",
  styleUrls: ["./test-suites.component.css"]
})
export class TestSuitesComponent implements OnInit {
  testSuitePath = undefined;
  testCaseArray = [];
  testSuites = [];
  selectedTestCaseData = [];
  testSuiteData = {};
  testCases = [];
  testSuiteFileName = undefined;
  testSuiteFileNameTruncated = undefined;
  testSuiteSchema = {layout: true};
  constructor(
    private testSuiteService: TestSuiteService,
    private testCaseService: TestCaseService,
    private configurationService: ConfigurationService,
    private toastr: ToastsManager
  ) {}

  ngOnInit() {
    this.configurationService
      .getSchema$("assets/testSuite_schema.json")
      .subscribe((data: any) => {
        this.testSuiteSchema = data;
      });
  }

  onAdd() {
    this.testSuiteService.getAllTestSuite$(this.testSuitePath).subscribe(
      data => {
        this.testSuites = data;
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
    this.testCaseService
      .getAllCases$(this.testSuitePath)
      .subscribe((data: any) => (this.testCases = data));
  }

  onCancel() {
    this.toastr.success("Your data has been cleared", "Success!");
    this.testSuiteData = undefined;
    this.selectedTestCaseData = [];
    this.testSuiteFileName = undefined;
    this.onAdd();
  }

  truncateFileName() {
    if (this.testSuiteFileName) {
      this.testSuiteFileNameTruncated = this.testSuiteFileName.substring(
        0,
        this.testSuiteFileName.length - 4
      );
    }
  }

  onUpdate(formData) {
    this.truncateFileName();
    this.testCaseArray["testCases"] = this.selectedTestCaseData;
    Object.assign(formData, this.testCaseArray);
    this.testSuiteService
      .putTestSuite$(formData, this.testSuitePath, this.testSuiteFileNameTruncated)
      .subscribe(
        data => {
          this.toastr.success("Success!");
          this.onAdd();
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

  onSave(formData) {
    this.testCaseArray["testCases"] = this.selectedTestCaseData;
    Object.assign(formData, this.testCaseArray);
    this.testSuiteService.postTestSuite$(formData, this.testSuitePath).subscribe(
      data => {
        this.toastr.success("Your data has been saved!", "Success!");
        this.onAdd();
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

  onSelectSuite(testSuite, selectedIndex) {
    this.onAdd();
    this.testSuiteData = testSuite;
    this.selectedTestCaseData = testSuite.testCases;
    this.testSuiteFileName = testSuite.file;
  }
  updateSelected(testCase) {
    this.selectedTestCaseData.push(testCase);

  }
  onReverse(selectedIndex) {
    this.selectedTestCaseData.splice(selectedIndex, 1);
  }
}
