package resource

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/icsnju/apt-mesos/mesosproto"
)

func RangeUsedUpdate(taskRanges *mesosproto.Value_Ranges, nodeRanges *mesosproto.Value_Ranges) *mesosproto.Value_Ranges {
	var newNodeRanges []*mesosproto.Value_Range
	for _, taskRange := range taskRanges.GetRange() {
		for _, nodeRange := range nodeRanges.GetRange() {
			if RangeInside(nodeRange, taskRange) {
				subedLeftRange, subedRightRange := RangeSub(nodeRange, taskRange)
				if subedLeftRange != nil {
					newNodeRanges = append(newNodeRanges, subedLeftRange)
				}
				if subedRightRange != nil {
					newNodeRanges = append(newNodeRanges, subedRightRange)
				}
			}
		}
	}
	return &mesosproto.Value_Ranges{
		Range: newNodeRanges,
	}
}

func RangeInside(large *mesosproto.Value_Range, small *mesosproto.Value_Range) bool {
	return large.GetBegin() <= small.GetBegin() && large.GetEnd() >= small.GetEnd()
}

func RangeSub(large *mesosproto.Value_Range, small *mesosproto.Value_Range) (*mesosproto.Value_Range, *mesosproto.Value_Range) {
	subedLeftBegin := large.GetBegin()
	subedLeftEnd := small.GetBegin() - 1
	subedLeftRange := &mesosproto.Value_Range{
		Begin: &subedLeftBegin,
		End:   &subedLeftEnd,
	}
	if subedLeftEnd < subedLeftBegin {
		subedLeftRange = nil
	}

	subedRightBegin := small.GetEnd() + 1
	subedRightEnd := large.GetEnd()
	subedRightRange := &mesosproto.Value_Range{
		Begin: &subedRightBegin,
		End:   &subedRightEnd,
	}
	if subedRightBegin > subedRightEnd {
		subedRightRange = nil
	}
	return subedLeftRange, subedRightRange
}

func RangeAdd(rangeOne *mesosproto.Value_Ranges, rangeTwo *mesosproto.Value_Ranges) *mesosproto.Value_Ranges {
	line := make([]bool, 65536)
	for _, enumRange := range rangeOne.GetRange() {
		for index := enumRange.GetBegin(); index <= enumRange.GetEnd(); index++ {
			line[index] = true
		}
	}

	for _, enumRange := range rangeTwo.GetRange() {
		for index := enumRange.GetBegin(); index <= enumRange.GetEnd(); index++ {
			line[index] = true
		}
	}

	newRanges := &mesosproto.Value_Ranges{}
	for index := 0; index < len(line); {
		for index < len(line) && !line[index] {
			index++
		}
		if index < len(line) {
			uint64Begin := uint64(index)
			newRange := &mesosproto.Value_Range{
				Begin: &uint64Begin,
			}

			for index < len(line) && line[index] {
				index++
			}
			uint64End := uint64(index - 1)
			newRange.End = &uint64End
			newRanges.Range = append(newRanges.Range, newRange)
		}
	}

	return newRanges
}

func ParseRanges(text string) (*mesosproto.Value_Ranges, error) {
	reg := regexp.MustCompile(`([\d]+)`)
	ports := reg.FindAllString(text, -1)
	ranges := &mesosproto.Value_Ranges{}

	if len(ports)%2 != 0 {
		return nil, errors.New("Cannot parse range string: " + text)
	}

	for i := 0; i < len(ports)-1; i += 2 {
		i64Begin, err := strconv.Atoi(ports[i])
		if err != nil {
			return nil, errors.New("Cannot parse range string: " + text)
		}
		i64End, err := strconv.Atoi(ports[i+1])
		if err != nil {
			return nil, errors.New("Cannot parse range string: " + text)
		}

		ranges.Range = append(ranges.Range, &mesosproto.Value_Range{
			Begin: proto.Uint64(uint64(i64Begin)),
			End:   proto.Uint64(uint64(i64End)),
		})
	}
	return ranges, nil
}
