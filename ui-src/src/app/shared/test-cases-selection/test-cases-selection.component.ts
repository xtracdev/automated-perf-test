import { Component, OnInit, Input } from '@angular/core';
import {TestCaseService} from '../../test-cases/test-case.service'
@Component({
  selector: 'test-cases-selection',
  templateUrl: './test-cases-selection.component.html',
  styleUrls: ['./test-cases-selection.component.css']
})
export class TestCasesSelectionComponent implements OnInit {

  @Input() testSuiteDirPath;
  availableTestCases = [];
  allTestCases = this.availableTestCases;
  selectedTestCases = [];


  constructor(
    private testCaseService: TestCaseService
  ) {}

  ngOnInit() {
    this.testCaseService.getAllTestCases$(this.testSuiteDirPath).subscribe( (data) => {
      this.availableTestCases = data;
      this.allTestCases = data;
    });
  }

  selectOne(){

  }

  deSelectOne(){

  }

  selectAll(){

  }
  selectNone(){

  }

}
