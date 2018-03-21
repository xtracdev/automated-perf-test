import {Component, Input, Output, EventEmitter} from "@angular/core";

@Component({
  selector: "app-test-cases-list",
  templateUrl: "./test-cases-list.component.html",
  styleUrls: ["./test-cases-list.component.css"]
})
export class TestCasesListComponent {
  selectedTestCaseData = [];
  @Input() testCases;
  @Output() addToSelected = new EventEmitter();
  @Output() reverse = new EventEmitter();

  selectedCase(testCase, i) {
    this.selectedTestCaseData = testCase;
    this.addToSelected.emit(this.selectedTestCaseData);
    this.reverse.emit(i);
  }
}
