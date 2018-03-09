import { Component, OnInit, Input, Output, EventEmitter } from "@angular/core";
import { TestCaseService } from "../../test-cases/test-case.service";

@Component({
  selector: "app-test-cases-selection",
  templateUrl: "./test-cases-selection.component.html",
  styleUrls: ["./test-cases-selection.component.css"]
})
export class TestCasesSelectionComponent implements OnInit {
  @Input() testSuiteDirPath;
  @Input() availableTestCases;
  @Input() selectedTestCaseData = [];
  @Output() addToSelected = new EventEmitter();
  @Output() reverse = new EventEmitter();
  constructor(private testCaseService: TestCaseService) {}

  ngOnInit() {}
  onReverse(selectedIndex) {
    this.reverse.emit(selectedIndex);
  }
  onAddToSelected(testCase, selectedIndex) {
    this.addToSelected.emit(testCase);
  }
}
