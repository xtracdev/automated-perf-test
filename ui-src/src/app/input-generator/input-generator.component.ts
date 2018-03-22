import { Component, OnChanges, OnInit, Output, Input, SimpleChanges } from '@angular/core';

@Component({
  selector: 'app-input-generator',
  templateUrl: './input-generator.component.html',
  styleUrls: ['./input-generator.component.css']
})
export class InputGeneratorComponent implements OnInit, OnChanges {


  @Input() schema;
  @Output() usersData;

  configSchema = {};
  constructor() {
  }


  ngOnChanges(changes: SimpleChanges) {
    console.log(`ngOnChanges - data is ${this.schema}`);
    for (let key in changes) {  
      console.log(`${key} changed.
  Current: ${changes[key].currentValue}.
  Previous: ${changes[key].previousValue}`);
    }
  }

  ngOnInit() {

    this.configSchema = this.schema;
    console.log('HEY HEY COmputer', this.configSchema)
  }

}
