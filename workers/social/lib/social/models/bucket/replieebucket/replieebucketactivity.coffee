CActivity = require '../../activity'
CBucketActivity = require '../../activity/bucketactivity'

module.exports = class CReplieeBucketActivity extends CBucketActivity
  @share()
  @set
    encapsulatedBy  : CActivity
    schema          : CActivity.schema
    sharedMethods   : CActivity.sharedMethods
    sharedEvents    : CBucketActivity.sharedEvents
    relationships   : CBucketActivity.relationships