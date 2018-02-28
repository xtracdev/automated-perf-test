
import {
  Component,
  Input,
  Output,
  EventEmitter
} from '@angular/core';

@Component({
  selector: 'test-cases-list',
  templateUrl: './test-cases-list.component.html',
  styleUrls: ['./test-cases-list.component.css']
})


export class TestCasesListComponent {
  testCaseArray = [];
  selectedTestCaseData = [];
  @Input() testCases;
  @Output() t = new EventEmitter();


  selectedCase(testCase) {
    this.selectedTestCaseData.push(testCase);
    this.testCaseArray = this.selectedTestCaseData;
    console.log(this.testCaseArray)
    this.t.emit(this.testCaseArray);
    
  }
  changeValue() {
    console.log("clicked")
    this.t.emit("emit this value");
  }  

}
