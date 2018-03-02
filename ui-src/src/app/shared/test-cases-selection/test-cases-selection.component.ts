import { Component, OnInit, Input, Output, EventEmitter } from "@angular/core";
import { TestCaseService } from "../../test-cases/test-case.service";

@Component({
  selector: "test-cases-selection",
  templateUrl: "./test-cases-selection.component.html",
  styleUrls: ["./test-cases-selection.component.css"]
})
export class TestCasesSelectionComponent implements OnInit {
  @Input() testSuiteDirPath;
  @Input() availableTestCases;
  selectedTestCaseData= [];
  constructor(private testCaseService: TestCaseService) {}

  ngOnInit() {
    // this.testCaseService.getAllTestCases$(this.testSuiteDirPath).subscribe( (data) => {
    //   this.availableTestCases = data;
    // });
  }

  selectOne() {}

  deSelectOne() {}

  selectAll() {}
  selectNone() {}
  onReverse(i) {
    this.selectedTestCaseData.splice(i, 1);
  }
  onAddToSelected(e) {
    this.selectedTestCaseData = e;
  }
}
