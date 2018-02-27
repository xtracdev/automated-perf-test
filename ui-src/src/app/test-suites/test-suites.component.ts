import { Component, OnInit } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { TestSuitesServices } from "./test-suites.service";
import { ToastsManager } from "ng2-toastr/ng2-toastr";

@Component({
  selector: "app-test-suites",
  templateUrl: "./test-suites.component.html",
  styleUrls: ["./test-suites.component.css"]
})
export class TestSuitesComponent  implements OnInit{
  testSuitePath = undefined;
  xmlFileName = undefined;
  fileName = undefined;
  
  constructor(
    private testSuitesServices: TestSuitesServices,
    private toastr: ToastsManager,
    private http: HttpClient
  ) {}
 

    ngOnInit(){
    }
    
    fileSelector(event) {
      this.fileName = event.srcElement.files[0].name;
      this.xmlFileName = this.fileName;
    }  

    onDeleteFile(){
      this.testSuitesServices
        .deleteSuites$('C:/Users/a622123/go/src/github.com/xtracdev/automated-perf-test/config', this.xmlFileName)
        .subscribe(
          data => {
            this.fileName= data;
            this.toastr.success("File sucessfuly deleted!");
          }
        );
    }

}
