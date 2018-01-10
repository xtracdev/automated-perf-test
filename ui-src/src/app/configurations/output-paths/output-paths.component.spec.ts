import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OutputPathsComponent } from './output-paths.component';

describe('OutputPathsComponent', () => {
  let component: OutputPathsComponent;
  let fixture: ComponentFixture<OutputPathsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OutputPathsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OutputPathsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
