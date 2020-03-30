[![Build Status](https://github.com/fogonthedowns/diplomatic/workflows/Go/badge.svg)](https://github.com/fogonthedowns/diplomatic/actions)

0. Loop through moves and determine if the destination is contested

1. ✅ Resolve all uncontested moves
   - ✅ is the move possible? If so Update the move.ResolvedMove
   - ✅ if not? Hold 

2. ✅ Collect all support orders
   - ✅ determine if the support has been cut.

3. ✅ Collect all convoy orders
   - ✅ determine if the path has been disrupted
   - ✅ determine if the path is complete

4. resolve all contested moves
   - ✅ will the piece move out of the way?
   - ✅ is there support?

5. Edge cases
   - ✅ Spain
   - ✅ Bulgaria
   - ✅ St. Petersburg

6. Retreat phase:
  - ✅ Add dislodged concept to piece.
  - ✅ If a piece is dislodged and it is retreat phase, it must move or be destroyed. 
  - ✅ A unit can not retreat from the origin of the attack.
  - ✅ Non dislodged pieces may not move in retreat phase.
  - ✅ Update the territory.Owner by move.PieceOwner after fall retreat phase is adjudicated
  - ✅ make Phase a type rather than an int (see newPhase, phase and Phase)
  - when assigning edge case Territory Save, ensure master territory is assigned ownership

7. Build Phase:
  - ✅ Add concept of Victory Center, count points
  - ✅ Count active units
     - if active units exceeds Victory Centers a destroy must be issued
     - if active units is less than Victory Center count a build order can be issued - only on home centers
  - If a Destroy must be issued and one is not issued, destroy a random unit
  - Create a new unit on the map, if a successful build order is issued
  - Ensure the year is cycled at the end of the build phase

8. Victory:
  - if the centers exceeds 18 at the end of the year the player wins the game
  - a victory can be shared between players who agree to share a victory

9. Validations:
    - ✅ verify the piece exists
    - ✅ Verify you own the piece.
    - ✅ verify the piece is at the matching start location 

10. design considerations:
  - ✅ you can not attack your own territory
  - ✅ you must be able to issue support and convoy orders for others
  - ✅ Players who have not issued moves, will issue hold orders
  - Move accepts phase, should it? Should it be a lookup on game?
