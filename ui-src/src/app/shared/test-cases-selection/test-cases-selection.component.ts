import { Component, OnInit, Input, Output,  EventEmitter } from '@angular/core';
import {TestCaseService} from '../../test-cases/test-case.service'

@Component({
  selector: 'test-cases-selection',
  templateUrl: './test-cases-selection.component.html',
  styleUrls: ['./test-cases-selection.component.css']
})
export class TestCasesSelectionComponent implements OnInit {

  @Input() testSuiteDirPath;
 @Input() availableTestCases;
 @Input() t;
 @Output() t1 = new EventEmitter();
selectedTestCases = []
  constructor(
    private testCaseService: TestCaseService
  ) {}

  ngOnInit() {
    // this.testCaseService.getAllTestCases$(this.testSuiteDirPath).subscribe( (data) => {
    //   this.availableTestCases = data;

    // });
  }

  selectOne(){

  }

  deSelectOne(){

  }

  selectAll(){

  }
  selectNone(){

  }

  onT(e) {
    console.log("@@@")
    this.selectedTestCases = e;
    console.log(this.selectedTestCases)
    console.log(this.availableTestCases)
    console.log(e)
    
  }

}
