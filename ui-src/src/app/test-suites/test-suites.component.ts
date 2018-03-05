import { Component, OnInit, Input } from "@angular/core";
import { TestSuiteService } from "./test-suite.service";
import { TestCaseService } from "../test-cases/test-case.service";
import { ToastsManager } from "ng2-toastr/ng2-toastr";
import {TestCasesSelectionComponent} from "../shared/test-cases-selection/test-cases-selection.component"
const TEST_SUITE_PATH =
  "C:/Users/A586754/go/src/github.com/xtracdev/automated-perf-test/config";
@Component({
  selector: "app-test-suites",
  templateUrl: "./test-suites.component.html",
  styleUrls: ["./test-suites.component.css"]
})
export class TestSuitesComponent {
  testSuitePath = TEST_SUITE_PATH;
  testCaseArray = [];
  testSuites = [];
  selectedTestCaseData = [];
  testSuiteData = {};
  testCases = [];
  testSuiteFileName = undefined;
  testSuiteSchema = { layout: true };
  constructor(
    private testSuiteService: TestSuiteService,
    private testCaseService: TestCaseService,
    private toastr: ToastsManager
  ) {}

  ngOnInit() {
    this.testSuiteService
      .getSchema$("assets/testSuite_schema.json")
      .subscribe((data: any) => {
        this.testSuiteSchema = data;
      });
  }

  onAdd() {
    this.testSuiteService.getAllTestSuite$(TEST_SUITE_PATH).subscribe(
      data => {
        this.testSuites = data;
        console.log(this.testSuites)
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
    this.testCaseService
      .getAllTestCases$(TEST_SUITE_PATH)
      .subscribe((data: any) => (this.testCases = data));
      
  }

  onCancel() {
    //clear schema and moving info back into available (get method)
    // this.testSuiteService.getTestSuite$(TEST_SUITE_PATH, this.testSuiteFileName)
    //   .subscribe(
    //     data => {
    //       this.testSuiteData = data;
    //       this.toastr.success("Previous data reloaded!");
    //     },
    //     error => {
    //       this.testSuiteData = undefined;
    //     }
    //   );
  }
  onUpdate(data) {
    this.testCaseArray["testCases"] = this.selectedTestCaseData;
    Object.assign(data, this.testCaseArray);
    this.testSuiteService
      .putTestSuite$(this.testSuiteData, this.testSuitePath, this.testSuiteFileName)
      .subscribe(
        data => {
          this.toastr.success("Success!");
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
  

  onSave(data) {
   this.testCaseArray["testCases"] = this.selectedTestCaseData;
    Object.assign(data, this.testCaseArray);
    this.testSuiteService.postTestSuite$(data, TEST_SUITE_PATH).subscribe(
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
            console.log(data);
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

  onSelectAll() {
    this.testCaseService
      .getAllTestCases$(TEST_SUITE_PATH)
      .subscribe((data: any) => {
       // this.selectedTestCaseData = data;
        //this.testCaseArray["testCases"] = this.selectedTestCaseData;
      });
  }
  onSelectSuite(testSuite, i)
  {
    this.testSuiteData = testSuite;
    this.selectedTestCaseData = testSuite.testCases;
    this.testSuiteFileName = testSuite.name


  } 
  updateSelected(e){
    this.selectedTestCaseData = e;
    console.log(e);
  }
  onReverseOne() {}
  onReverseAll() {
  //  this.selectedTestCaseData = undefined;
  }
}
