# Classroom "sorting hat"

## inputs

 * list of pupils
 * list of siblings already at school
 * classes (three houses, 2 or 3 classes per house, e.g. 1a1, 1a2, 1b1, 1b2, 1c1, 1c2, 1c3)

## rules

 1. pupils with siblings in the school should be in the same class (e.g. 2a2 -> 1a2), or if not possible, at least the same house
 1. pupils can specify other desired pupils to be in a class with
 1. pupils can list others they do not want to be in a class with
 1. class sizes should be as even as possible and not exceeding 30
 1. there should be an even spread of pupils with GIRFME plans

## input checking

There may be conflicting rules, e.g. 

 * pupil A wants to be with B and C, but C does not want to be with A.
 * or pupil D and E have siblings placing them in 1a1 but D does not want to be in a class with E

## approaches

### input validation and correction

Build a graph of pupil (dis-)assocation requests.  pupils are nodes, edges will be "associate" and "disassociate".  Presenting these input graphs visually will likely help users understand the domain better (look at https://github.com/go-echarts/go-echarts).

 * Ideally, this will form a list of disjoint graphs each of which has < 30 nodes and contains no cycles of both edge types
 * where the graphs are too large, we need to decide how to split them
 * where there are cycles containing differing edge types, we need to decide which takes preference
  * e.g. A wants to be with B, C; C wants to be with D, E; E does not want to be with A

we should adjust the rules such that they are possible to satisfy prior to trying to place pupils.  This could be done automatically based on priority order, or presented back to the user to perform manual correction

### seeding

The pupils with siblings should be placed first as this very simple and usually not very numerous

### Monte Carlo method

We can randomly assign pupils to classes, then score each class on how well it satisfies the rules.  after a sufficient number of attempts have been made we should have approached an optimal solution

The most optimal solution should be presented with a list of rules that were unable to be satisfied (if any).  The user can then complete any final tweaks manually

## implementation

Performance is key as many millions of simulations will need to be tested given approx 200 pupils and 7 classes.  Datastructures should be chosen to be as optimal as possible for each task. As the problem space is not huge, duplication is fine if different stuctures of the same data will improve separate steps.  Obviously algorithms should be as optimal as possible working in tandem with the chosen structures.  we should use multithreading liberally.  initially targeting a 6c/12t machine, we should be able to configure optimal parallelism.

### multithreading

### datastructures
