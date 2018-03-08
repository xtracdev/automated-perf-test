import { Component, OnInit } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { ConfigurationService } from "./configuration.service";
import { TestSuiteService } from "../test-suites/test-suite.service";
import { JsonSchemaFormModule } from "angular2-json-schema-form";
import { ToastsManager } from "ng2-toastr/ng2-toastr";
import "rxjs/add/operator/map";

@Component({
  selector: "app-configurations",
  templateUrl: "./configurations.component.html",
  styleUrls: ["./configurations.component.css"]
})
export class ConfigurationsComponent implements OnInit {
  formData = {};
  test = [];
  test2 = [];
  testArary =[];
  configPath = undefined;
  xmlFileName = undefined;
  fileName = undefined;
  // needed for layout to load
  configSchema = {layout:true};

  constructor(
    private configurationService: ConfigurationService,
    private testSuiteService: TestSuiteService,
    private toastr: ToastsManager,

    private http: HttpClient
  ) {}

  ngOnInit() {
    this.configurationService
      .getSchema$("assets/schema.json")
      .subscribe((data: any) => {
        this.configSchema = data;
        console.log("here",this.configSchema)
      });
    //   this.testSuiteService.getAllTestSuite$("C:/Users/A586754/go/src/github.com/xtracdev/automated-perf-test/config").subscribe(
    //     data => {
    //       this.test = data
    //       for (var i = 0; i < this.test.length; i++ ) {
    //        this.testArary.push(this.test[i].name)
            
            
    //         console.log("helloo", this.configSchema)
    //       }
    // this.configSchema["properties"].testSuite.enum = this.testArary
    // this.test2 = this.configSchema["properties"].testSuite.enum
    // this.formData["testSuite"] = this.test2
    //       console.log("bye",this.formData)
    //       console.log(this.testArary)
    //       console.log("configSchema", this.configSchema)
    //       //this.toastr.success("Success!");
    //     },
    
    //     error => {
    //       switch (error.status) {
    //         case 500: {
    //           this.toastr.error("An error has occurred!", "Check the logs!");
    //           break;
    //         }
    //         case 400: {
    //           this.toastr.error(
    //             "No Test Suite Directory added",
    //             "An error occurred!"
    //           );
    //           break;
    //         }
    //         default: {
    //           this.toastr.error("An error occurred!");
    //         }
    //       }
    //     }
    //   );
    }   
  
  onGet(){
  this.testSuiteService.getAllTestSuite$("C:/Users/A586754/go/src/github.com/xtracdev/automated-perf-test/config").subscribe(
    data => {
      this.test = data
      for (var i = 0; i < this.test.length; i++ ) {
       this.testArary.push(this.test[i].name)
        
        
        console.log("helloo", this.configSchema)
      }
      this.test2 = this.testArary;
this.configSchema["properties"].testSuite.enum = this.test2;
      console.log(this.testArary)
      console.log("configSchema", this.configSchema)
      //this.toastr.success("Success!");
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

  onSubmit(configData) {

    this.configurationService.postConfig$(configData, this.configPath).subscribe(
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

  onCancel() {
    this.configurationService
      .getConfig$(this.configPath, this.xmlFileName)
      .subscribe(
        data => {
          this.formData = data;
          this.toastr.success("Previous data reloaded!");
        },
        error => {
          this.configPath = undefined;
          this.xmlFileName = undefined;
          this.formData = undefined;
        }
      );
  }

  fileSelector(event) {
    this.fileName = event.srcElement.files[0].name;
    this.xmlFileName = this.fileName;
    this.xmlFileName = this.xmlFileName.substring(
      0,
      this.xmlFileName.length - 4);


    this.onGetFile();
    (<HTMLInputElement>document.getElementById("file")).value ="";
   
  }

  onClearFile(){
    this.xmlFileName = undefined;
  }
  onGetFile() {
    this.configurationService
      .getConfig$(this.configPath, this.xmlFileName)
      .subscribe(
        data => {
          this.formData = data;
          console.log(this.formData)
          this.toastr.success("Success!");
        },
        error => {
          switch (error.status) {
            case 404: {
              this.toastr.error("File not found!", "An error occured!");
              break;
            }
            case 400: {
              this.toastr.error(
                "Check your field inputs" ,
                "An error occurred!"
              );
              break;
            }
            case 500: {
              this.toastr.error("An error has occurred!", "Check the logs!");
              break;
            }
            default: {
              this.toastr.error(
                "Your data was not retrieved!",
                "An error occurred!"
              );
            }
          }
        }
      );
  }
  onUpdate(configData) {
    this.configurationService
      .putConfig$(this.formData, this.configPath, this.xmlFileName)
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
}
